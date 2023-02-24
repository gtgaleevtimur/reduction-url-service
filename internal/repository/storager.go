package repository

import (
	"context"
	"errors"
)

// Storager - интерфейс хранилища.
type Storager interface {
	GetShortURL(ctx context.Context, fullURL string) (string, error)
	GetFullURL(ctx context.Context, shortURL string) (string, error)
	saveData(ctx context.Context, fullURL string, userid string, hash string) error
	InsertURL(ctx context.Context, fURL string, userID string) (string, error)
	GetAllUserURLs(ctx context.Context, userid string) ([]SlicedURL, error)
	Delete(ctx context.Context, hashes []string, userID string) error
	Ping(ctx context.Context) error
	GetCountURL(ctx context.Context) (int, error)
	GetCountUsers(ctx context.Context) (int, error)
}

// NodeURL - сущность сокращенного URL, использующаяся в логике резервного хранилища.
type NodeURL struct {
	Hash   string `json:"hash"`
	FURL   string `json:"original_url"`
	UserID string `json:"user_id"`
	Delete bool   `json:"is_deleted"`
}

// URL - сущность URL, использующаяся для записи в хэш-таблице по hash-ключу сокращенного URL.
type URL struct {
	UserID string `json:"userid"`
	FURL   string `json:"original_url"`
	Delete bool   `json:"is_deleted"`
}

// FullURL - сущность URL, использующая для записи оригинального URL в эндпоинта POST /api/shorten принимающего JSON.
type FullURL struct {
	Full string `json:"url"`
}

// ShortURL - сущность URL, использующая для ответа сокращенного URL в эндпоинта POST /api/shorten принимающего JSON.
type ShortURL struct {
	Short string `json:"result"`
}

// SlicedURL - сущность URL, использующаяся для формирования ответа с массивом всех сохраненных пользователем URL.
type SlicedURL struct {
	Short string `json:"short_url" db:"hash"`
	Full  string `json:"original_url" db:"url"`
}

// FullBatch - сущность URL, использующаяся для записи массива с URL в эндпоинте POST /api/shorten/batch.
type FullBatch struct {
	CorID string `json:"correlation_id"`
	Full  string `json:"original_url"`
}

// ShortBatch - сущность URL, использующаяся для ответа  в эндпоинте POST /api/shorten/batch.
type ShortBatch struct {
	CorID string `json:"correlation_id"`
	Short string `json:"short_url"`
}

// StatStruct - сущность статистики сокращенных URL и количества пользователей.
type StatStruct struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}

// ErrConflictInsert - ошибка, показывающая, что сохраняемый URL уже есть в базе данных.
var ErrConflictInsert error = errors.New("URL is exist")

// ErrNotFoundURL - ошибка,показывающая , что запрашиваемый URL нет в базе данных.
var ErrNotFoundURL error = errors.New("URL not found in DB")

// ErrDeletedURL - ошибка,показывающая , что запрашиваемый URL нет удален из БД.
var ErrDeletedURL error = errors.New("URL is delete")

// ErrFileStoragePathNil - ошибка, показывающая, что путь записи резервного хранилища не задан.
var ErrFileStoragePathNil error = errors.New("err FILE_STORAGE_PATH is nil ")
