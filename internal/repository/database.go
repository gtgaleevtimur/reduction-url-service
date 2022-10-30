package repository

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
)

type Database struct {
	DB *sql.DB
}

func NewDatabaseDSN(conf *config.Config) (Storager, error) {
	db, err := sql.Open("pgx", conf.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	table := `CREATE TABLE IF NOT EXISTS "shortener" ("hash" TEXT UNIQUE PRIMARY KEY NOT NULL,"url" TEXT UNIQUE NOT NULL,"userid" TEXT NOT NULL)`
	_, err = db.Exec(table)
	if err != nil {
		return nil, err
	}

	s := &Database{
		DB: db,
	}
	return s, nil
}

func (d *Database) GetShortURL(fullURL string) (string, error) {
	var hash string
	str := `SELECT "hash" FROM "shortener" WHERE "url" = $1`
	err := d.DB.QueryRow(str, fullURL).Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (d *Database) GetFullURL(shortURL string) (string, error) {
	var fullURL string
	str := `SELECT "url" FROM "shortener" WHERE "hash" = $1`
	err := d.DB.QueryRow(str, shortURL).Scan(&fullURL)
	if err != nil {
		return "", err
	}
	return fullURL, nil
}

func (d *Database) InsertURL(fullURL string, userid string, hash string) error {
	if fullURL == "" || fullURL == " " || userid == "" || userid == " " || hash == "" || hash == " " {
		return errors.New("ErrNoEmptyInsert")
	}
	tr, err := d.DB.Begin()
	if err != nil {
		return err
	}
	defer tr.Rollback()
	str := `INSERT INTO "shortener"("hash","url","userid")VALUES ($1,$2,$3)`
	st, err := d.DB.Prepare(str)
	if err != nil {
		return err
	}
	defer st.Close()
	_, err = st.Exec(hash, fullURL, userid)
	if err != nil {
		return err
	}
	err = tr.Commit()
	if err != nil {
		return err
	}
	return nil
}

//MiddlewareInsert - метод-помощник, генерирует hash для ключа,передает hash+url+userid хранилищу,возвращает сокращенный url
func (d *Database) MiddlewareInsert(fURL string, userID string) (string, error) {
	//Генерируем hash.
	hasher := md5.Sum([]byte(fURL + userID))
	hash := hex.EncodeToString(hasher[:len(hasher)/5])
	//Проверяем есть ли в хранилище такой url.
	okHash, err := d.GetShortURL(fURL)
	//Если нет,то вставляем новые данные.
	if err != nil {
		err = d.InsertURL(fURL, userID, hash)
		if err != nil {
			return "", err
		}
		//Возвращаем сгенерированный hash.
		return hash, nil
	}
	//Если есть , возвращаем hash.
	return okHash, nil
}

func (d *Database) GetAllUserURLs(userid string) ([]SlicedURL, error) {
	var hash string
	var fullURL string
	result := make([]SlicedURL, 0)

	str := `SELECT "hash", "url" FROM "shortener" WHERE "userid" = $1`
	rows, err := d.DB.Query(str, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&hash, &fullURL)
		if err != nil {
			return nil, err
		}
		result = append(result, SlicedURL{
			Short: hash,
			Full:  fullURL,
		})
	}
	return result, nil
}

func (d *Database) Ping() error {
	return d.DB.Ping()
}
