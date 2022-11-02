package repository

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
)

//Database - структура базы данных SQL.
type Database struct {
	DB *sql.DB
}

//NewDatabaseDSN - конструктор базы данных на основе SQL,возвращает интерфейс.
func NewDatabaseDSN(conf *config.Config) (Storager, error) {
	//Инициализируем драйвер SQL.
	db, err := sql.Open("pgx", conf.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	//Проверка Ping.
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	//Подготавливаем SQL запрос на создание таблицы,если ее нет.
	table := `CREATE TABLE IF NOT EXISTS "shortener" ("hash" TEXT UNIQUE PRIMARY KEY NOT NULL,"url" TEXT UNIQUE NOT NULL,"userid" TEXT NOT NULL)`
	//Выполняем SQL запрос.
	_, err = db.Exec(table)
	if err != nil {
		return nil, err
	}
	//Возвращаем интерфейс.
	s := &Database{
		DB: db,
	}
	return s, nil
}

//GetShortURL - метод, возвращающий hash сокращенного url.
func (d *Database) GetShortURL(fullURL string) (string, error) {
	var hash string
	//Готовим SQL запрос и выполняем.
	str := `SELECT "hash" FROM "shortener" WHERE "url" = $1`
	err := d.DB.QueryRow(str, fullURL).Scan(&hash)
	if err != nil {
		return "", err
	}
	//Возвращем hash,если не было ошибки.
	return hash, nil
}

//GetFullURL - метод, возвращающий original_url по его hash.
func (d *Database) GetFullURL(hash string) (string, error) {
	var fullURL string
	//Готовим SQL запрос и выполняем.
	str := `SELECT "url" FROM "shortener" WHERE "hash" = $1`
	err := d.DB.QueryRow(str, hash).Scan(&fullURL)
	if err != nil {
		return "", err
	}
	//Возвращем original_url,если не было ошибки.
	return fullURL, nil
}

//InsertURL - метод, который сохраняет original_url,user_id и hash в базу данных.
func (d *Database) InsertURL(fullURL string, userid string, hash string) error {
	//Проверяем полученные данные.
	if fullURL == "" || fullURL == " " || userid == "" || userid == " " || hash == "" || hash == " " {
		return errors.New("ErrNoEmptyInsert")
	}
	//Объявляем начало транзакции.
	tr, err := d.DB.Begin()
	if err != nil {
		return err
	}
	defer tr.Rollback()
	//Подготавливаем стейтмент для БД.
	str := `INSERT INTO "shortener"("hash","url","userid")VALUES ($1,$2,$3)`
	st, err := d.DB.Prepare(str)
	if err != nil {
		return err
	}
	defer st.Close()
	//Выполняем стейтмент.
	_, err = st.Exec(hash, fullURL, userid)
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

//MiddlewareInsert - метод-помощник, генерирует hash для ключа,передает hash+url+userid хранилищу,возвращает сокращенный url
func (d *Database) MiddlewareInsert(fullURL string, userID string) (string, error) {
	//Генерируем hash.
	hasher := md5.Sum([]byte(fullURL + userID))
	hash := hex.EncodeToString(hasher[:len(hasher)/5])
	//Проверяем есть ли в хранилище такой url.
	okHash, err := d.GetShortURL(fullURL)
	//Если нет,то вставляем новые данные.
	if err != nil {
		err = d.InsertURL(fullURL, userID, hash)
		if err != nil {
			return "", err
		}
		//Возвращаем сгенерированный hash.
		return hash, nil
	}
	//Если есть , возвращаем hash и ошибку.
	return okHash, ErrConflictInsert
}

//GetAllUserURLs - метод возвращающий массив со всеми original_url+hash сохраненными пользователем.
func (d *Database) GetAllUserURLs(userid string) ([]SlicedURL, error) {
	//Объявляем переменные и массив с результатом.
	var hash string
	var url string
	result := make([]SlicedURL, 0)
	//Подготавливаем/выполняем запрос базе данных.
	str := `SELECT "hash", "url" FROM "shortener" WHERE "userid" = $1`
	rows, err := d.DB.Query(str, userid)
	//Проверяем обе! ошибки.
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

//Ping - возвращает ответ от БД Ping.
func (d *Database) Ping() error {
	return d.DB.Ping()
}
