package repository

import (
	"database/sql"

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
	} else {
		table := `CREATE TABLE IF NOT EXIST "shortener" ("hash" TEXT UNIQUE PRIMARY KEY NOT NULL,"url" TEXT UNIQUE NOT NULL,"userid" TEXT NOT NULL)`
		_, err = db.Exec(table)
		if err != nil {
			return nil, err
		}
	}
	s := &Database{
		DB: db,
	}
	return s, nil
}

func (d *Database) GetShortURL(fullURL string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Database) GetFullURL(shortURL string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Database) InsertURL(fullURL string, userid string, hash string) error {
	//TODO implement me
	panic("implement me")
}

func (d *Database) LoadRecoveryStorage(str string) error {
	//TODO implement me
	panic("implement me")
}

func (d *Database) MiddlewareInsert(fURL string, userID string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Database) GetAllUserURLs(userid string) ([]SlicedURL, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Database) Ping() error {
	return d.DB.Ping()
}
