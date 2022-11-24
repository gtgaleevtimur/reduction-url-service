package repository

import (
	"context"
	"errors"
)

// Storager - интерфейс хранилища.
type Storager interface {
	GetShortURL(ctx context.Context, fullURL string) (string, error)
	GetFullURL(ctx context.Context, shortURL string) (string, int)
	saveData(ctx context.Context, fullURL string, userid string, hash string) error
	InsertURL(ctx context.Context, fURL string, userID string) (string, error)
	GetAllUserURLs(ctx context.Context, userid string) ([]SlicedURL, error)
	Delete(userID []string, shortURL string) error
	Ping(ctx context.Context) error
}
type Task struct {
	UserID string
	Hash   string
}

type NodeURL struct {
	Hash   string `json:"hash"`
	FURL   string `json:"original_url"`
	UserID string `json:"user_id"`
	Delete bool   `json:"is_deleted"`
}

type URL struct {
	UserID string `json:"userid"`
	FURL   string `json:"original_url"`
	Delete bool   `json:"is_deleted"`
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

// ErrConflictInsert - ошибка, показывающая, что сохраняемый URL уже есть в базе данных.
var ErrConflictInsert error = errors.New("URL is exist")

// ErrNotFoundURL - ошибка,показывающая , что запрашиваемый URL нет в базе данных.
var ErrNotFoundURL error = errors.New("URL not found in DB")

// ErrDeletedURL - ошибка,показывающая , что запрашиваемый URL нет удален из БД.
var ErrDeletedURL error = errors.New("URL is delete")
