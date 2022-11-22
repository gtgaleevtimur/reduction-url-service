package repository

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
)

// Database - структура базы данных SQL.
type Database struct {
	DB *sql.DB
}

// NewDatabaseDSN - конструктор базы данных на основе SQL, возвращает интерфейс.
func NewDatabaseDSN(conf *config.Config) (Storager, error) {
	s := &Database{}
	err := s.Connect(conf)
	if err != nil {
		return nil, err
	}
	err = s.Bootstrap()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (d *Database) Bootstrap() (err error) {
	//Подготавливаем SQL запрос на создание таблицы, если ее нет.
	table := `CREATE TABLE IF NOT EXISTS "shortener" ("hash" TEXT UNIQUE PRIMARY KEY NOT NULL,
													"url" TEXT UNIQUE NOT NULL,
													"userid" TEXT NOT NULL,
													"delete" BOOLEAN NOT NULL)`
	//Выполняем SQL запрос.
	_, err = d.DB.Exec(table)
	if err != nil {
		return err
	}
	return nil
}

// Connect - метод выполняет соединение с базой данных.
func (d *Database) Connect(conf *config.Config) (err error) {
	d.DB, err = sql.Open("pgx", conf.DatabaseDSN)
	if err != nil {
		return err
	}
	err = d.DB.Ping()
	if err != nil {
		return err
	}
	return nil
}

// GetShortURL - метод, возвращающий hash сокращенного url.
func (d *Database) GetShortURL(ctx context.Context, fullURL string) (string, error) {
	var hash string
	//Задаем контекст на основе переданного из запроса
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	//Готовим SQL запрос и выполняем.
	str := `SELECT "hash" FROM "shortener" WHERE "url" = ($1) AND "delete" = false`
	err := d.DB.QueryRowContext(ctx, str, fullURL).Scan(&hash)
	if err != nil {
		return "", err
	}
	//Возвращем hash, если не было ошибки.
	return hash, nil
}

// GetFullURL - метод, возвращающий original_url по его hash.
func (d *Database) GetFullURL(ctx context.Context, hash string) (string, error) {
	var fullURL string
	var del bool
	//Задаем контекст на основе переданного из запроса
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	//Готовим SQL запрос и выполняем.
	str := `SELECT "url", "delete" FROM "shortener" WHERE "hash" = ($1)`
	err := d.DB.QueryRowContext(ctx, str, hash).Scan(&fullURL, &del)
	if err != nil {
		return "", err
	}
	//Если URL удален возвращаем соответствующую ошибку.
	if del {
		return "", ErrDeletedURL
	}
	//Возвращем original_url, если не было ошибки.
	return fullURL, nil
}

// InsertURL - метод, который сохраняет original_url,user_id и hash в базу данных.
func (d *Database) saveData(ctx context.Context, fullURL string, userid string, hash string) error {
	//Проверяем полученные данные.
	if fullURL == "" || fullURL == " " || userid == "" || userid == " " || hash == "" || hash == " " {
		return errors.New("ErrNoEmptyInsert")
	}
	//Задаем контекст на основе переданного из запроса
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	//Объявляем начало транзакции.
	tr, err := d.DB.Begin()
	if err != nil {
		return err
	}
	//defer tr.Rollback()
	//Подготавливаем стейтмент для БД.
	str := `INSERT INTO "shortener"("hash","url","userid","delete")VALUES ($1,$2,$3,false)`
	st, err := tr.Prepare(str)
	if err != nil {
		return err
	}
	defer st.Close()
	//Выполняем стейтмент.
	_, err = st.ExecContext(ctx, hash, fullURL, userid)
	if err != nil {
		return err
	}
	//Подтверждаем транзакцию.
	err = tr.Commit()
	if err != nil {
		return err
	}
	return nil
}

// InsertURL - метод ,который генерирует hash для ключа,передает hash+url+userid хранилищу,возвращает сокращенный url
func (d *Database) InsertURL(ctx context.Context, fullURL string, userID string) (string, error) {
	//Генерируем hash.
	hasher := md5.Sum([]byte(fullURL + userID))
	hash := hex.EncodeToString(hasher[:len(hasher)/5])
	//Проверяем есть ли в хранилище такой url.
	okHash, err := d.GetShortURL(ctx, fullURL)
	//Если нет, то вставляем новые данные.
	if err != nil {
		err = d.saveData(ctx, fullURL, userID, hash)
		if err != nil {
			return "", err
		}
		//Возвращаем сгенерированный hash.
		return hash, nil
	}
	//Если есть, возвращаем hash и ошибку.
	return okHash, ErrConflictInsert
}

// GetAllUserURLs - метод возвращающий массив со всеми original_url+hash сохраненными пользователем.
func (d *Database) GetAllUserURLs(ctx context.Context, userid string) ([]SlicedURL, error) {
	//Задаем контекст на основе переданного из запроса
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	//Объявляем переменные и массив с результатом.
	var hash string
	var url string
	result := make([]SlicedURL, 0)
	//Подготавливаем/выполняем запрос базе данных.
	str := `SELECT "hash", "url" FROM "shortener" WHERE "userid" = ($1) AND "delete" = false`
	rows, err := d.DB.QueryContext(ctx, str, userid)
	//Проверяем обе ошибки.
	if err != nil || rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()
	//Итерируемся внутри полученного курсора.
	for rows.Next() {
		//Сканируем строку в переменные.
		err = rows.Scan(&hash, &url)
		if err != nil {
			return nil, err
		}
		//Заполняем массив с результатом.
		result = append(result, SlicedURL{
			Short: hash,
			Full:  url,
		})
	}
	return result, nil
}

// Ping - возвращает ответ от БД Ping.
func (d *Database) Ping(ctx context.Context) error {
	//Задаем контекст на основе переданного из запроса
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return d.DB.PingContext(ctx)
}

// Delete - метод, который данные помечает как удаленные по их hash(идентификатор).
func (d *Database) Delete(ctx context.Context, shortURL string, userID string) error {
	//Объявляем начало транзакции.
	tr, err := d.DB.Begin()
	if err != nil {
		return err
	}
	//	defer tr.Rollback()
	//Подготавливаем стейтмент для БД.
	str := `UPDATE "shortener" SET "delete" = true WHERE "hash" = ($1) and "userid" = ($2)`
	st, err := tr.Prepare(str)
	if err != nil {
		return err
	}
	defer st.Close()
	//Выполняем транзакцию, передавая драйверу массив с идентификаторами URL.
	if _, err = st.ExecContext(ctx, shortURL, userID); err != nil {
		return err
	}
	//Возвращаем результат транзакции.
	return tr.Commit()
}
