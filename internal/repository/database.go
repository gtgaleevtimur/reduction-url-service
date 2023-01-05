package repository

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
)

// Database - структура базы данных SQL.
type Database struct {
	DB *sql.DB
	sync.Mutex
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

// Bootstrap - метод, создающий рабочую таблицу в БД.
func (d *Database) Bootstrap() (err error) {
	//Подготавливаем SQL запрос на создание таблицы, если ее нет.
	// Выполняем SQL запрос.
	_, err = d.DB.Exec(`CREATE TABLE IF NOT EXISTS shortener (hashid TEXT UNIQUE PRIMARY KEY NOT NULL,
													url TEXT UNIQUE NOT NULL,
													userid TEXT NOT NULL,
													is_deleted BOOLEAN NOT NULL)`)
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
	// Готовим SQL запрос и выполняем.
	err := d.DB.QueryRowContext(ctx, `SELECT hashid FROM shortener WHERE url = $1 AND is_deleted = false`, fullURL).Scan(&hash)
	if err != nil {
		return "", err
	}
	// Возвращем hash, если не было ошибки.
	return hash, nil
}

// GetFullURL - метод, возвращающий original_url по его hash.
func (d *Database) GetFullURL(ctx context.Context, hash string) (string, error) {
	var fullURL string
	var del bool
	// Готовим SQL запрос и выполняем.
	err := d.DB.QueryRowContext(ctx, `SELECT url , is_deleted FROM shortener WHERE hashid = $1`, hash).Scan(&fullURL, &del)
	// Если URL отсутствует в БД возвращаем соответствующую ошибку.
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFoundURL
	}
	// Если URL удален возвращаем соответствующую ошибку.
	if del {
		return "", ErrDeletedURL
	}
	// Возвращем original_url.
	return fullURL, nil
}

// InsertURL - метод, который сохраняет original_url,user_id и hash в базу данных.
func (d *Database) saveData(ctx context.Context, fullURL string, userid string, hash string) error {
	// Проверяем полученные данные.
	if fullURL == "" || fullURL == " " || userid == "" || userid == " " || hash == "" || hash == " " {
		return errors.New("ErrNoEmptyInsert")
	}
	// Объявляем начало транзакции.
	tr, err := d.DB.Begin()
	if err != nil {
		return err
	}
	defer tr.Rollback()
	// Подготавливаем стейтмент для БД.
	st, err := tr.Prepare(`INSERT INTO shortener(hashid,url,userid,is_deleted)VALUES ($1,$2,$3,false)`)
	if err != nil {
		return err
	}
	defer st.Close()
	// Выполняем стейтмент.
	_, err = st.ExecContext(ctx, hash, fullURL, userid)
	if err != nil {
		return err
	}
	// Подтверждаем транзакцию.
	err = tr.Commit()
	if err != nil {
		return err
	}
	return nil
}

// InsertURL - метод ,который генерирует hash для ключа,передает hash+url+userid хранилищу,возвращает сокращенный url.
func (d *Database) InsertURL(ctx context.Context, fullURL string, userID string) (string, error) {
	// Генерируем hash.
	hasher := md5.Sum([]byte(fullURL + userID))
	hash := hex.EncodeToString(hasher[:len(hasher)/5])
	// Проверяем есть ли в хранилище такой url.
	okHash, err := d.GetShortURL(ctx, fullURL)
	// Если нет, то вставляем новые данные.
	if err != nil {
		err = d.saveData(ctx, fullURL, userID, hash)
		if err != nil {
			return "", err
		}
		// Возвращаем сгенерированный hash.
		return hash, nil
	}
	// Если есть, возвращаем hash и ошибку.
	return okHash, ErrConflictInsert
}

// GetAllUserURLs - метод возвращающий массив со всеми original_url+hash сохраненными пользователем.
func (d *Database) GetAllUserURLs(ctx context.Context, userid string) ([]SlicedURL, error) {
	// Объявляем переменные и массив с результатом.
	var hash string
	var url string
	result := make([]SlicedURL, 0)
	// Подготавливаем/выполняем запрос базе данных.
	rows, err := d.DB.QueryContext(ctx, `SELECT hashid , url FROM shortener WHERE userid = $1 AND is_deleted = false`, userid)
	// Проверяем обе ошибки.
	if err != nil || rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()
	// Итерируемся внутри полученного курсора.
	for rows.Next() {
		// Сканируем строку в переменные.
		err = rows.Scan(&hash, &url)
		if err != nil {
			return nil, err
		}
		// Заполняем массив с результатом.
		result = append(result, SlicedURL{
			Short: hash,
			Full:  url,
		})
	}
	return result, nil
}

// Ping - возвращает ответ от БД Ping.
func (d *Database) Ping(ctx context.Context) error {
	// Задаем контекст на основе переданного из запроса.
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return d.DB.PingContext(ctx)
}

// Delete - метод, который данные помечает как удаленные по их hash(идентификатор).
func (d *Database) Delete(ctx context.Context, hashes []string, userID string) error {
	d.Lock()
	defer d.Unlock()
	// Инициализируем контекст с таймаутом.
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()
	// Объявляем начало транзакции.
	tr, err := d.DB.Begin()
	if err != nil {
		return err
	}
	defer tr.Rollback()
	// Подготавливаем стейтмент для БД.
	st, err := tr.Prepare(`update shortener set is_deleted=true WHERE hashid = any ($1) and userid = $2`)
	if err != nil {
		return err
	}
	defer st.Close()
	// Выполняем транзакцию.
	if _, err = st.ExecContext(ctx, hashes, userID); err != nil {
		return err
	}
	// Возвращаем результат транзакции.
	return tr.Commit()
}

// clearTable - хелпер-метод, очищающий поля таблицы.
func (d *Database) clearTable() error {
	_, err := d.DB.Exec(`delete from shortener`)
	return err
}
