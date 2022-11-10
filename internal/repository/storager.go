package repository

import (
	"errors"
)

// Storager - интерфейс хранилища.
type Storager interface {
	GetShortURL(fullURL string) (string, error)
	GetFullURL(shortURL string) (string, error)
	saveData(fullURL string, userid string, hash string) error
	InsertURL(fURL string, userID string) (string, error)
	GetAllUserURLs(userid string) ([]SlicedURL, error)
	Ping() error
}

type NodeURL struct {
	Hash   string `json:"hash"`
	FURL   string `json:"original_url"`
	UserID string `json:"user_id"`
}

type URL struct {
	UserID string `json:"userid"`
	FURL   string `json:"original_url"`
}

type FullURL struct {
	Full string `json:"url"`
}

type ShortURL struct {
	Short string `json:"result"`
}

type SlicedURL struct {
	Short string `json:"short_url" db:"hash"`
	Full  string `json:"original_url" db:"url"`
}

type FullBatch struct {
	CorID string `json:"correlation_id"`
	Full  string `json:"original_url"`
}

type ShortBatch struct {
	CorID string `json:"correlation_id"`
	Short string `json:"short_url"`
}

//ErrConflictInsert - ошибка,показывающая что сохраняемый URL уже есть в базе данных.
var ErrConflictInsert error = errors.New("URL is exist")
