package repository

import (
	"context"
	"errors"
	"strconv"
	"sync"
)

type FullURL struct {
	Full string `json:"url"`
}

type ShortURL struct {
	Short string `json:"result"`
}

type Storage struct {
	Counter        int
	FullURLKeyMap  map[string]ShortURL
	ShortURLKeyMap map[string]FullURL
	sync.Mutex
}

func NewStorage() *Storage {
	s := &Storage{
		Counter:        0,
		FullURLKeyMap:  make(map[string]ShortURL),
		ShortURLKeyMap: make(map[string]FullURL),
	}
	return s
}

func (s *Storage) GetShortURL(_ context.Context, fullURL string) (string, error) {
	s.Lock()
	defer s.Unlock()
	if val, ok := s.FullURLKeyMap[fullURL]; ok {
		return val.Short, nil
	}
	return "", errors.New("wrong URL")
}

func (s *Storage) GetFullURL(_ context.Context, shortURL string) (string, error) {
	s.Lock()
	defer s.Unlock()
	if val, ok := s.ShortURLKeyMap[shortURL]; ok {
		return val.Full, nil
	}
	return "", errors.New("wrong URL")
}

func (s *Storage) InsertURL(ctx context.Context, fullURL string) (string, error) {
	if fullURL == "" || fullURL == " " {
		return "", errors.New("ErrNoNilInsert")
	}
	short, err := s.GetShortURL(ctx, fullURL)
	if err == nil {
		return short, nil
	}
	s.Lock()
	defer s.Unlock()
	fURL := FullURL{Full: fullURL}
	sURL := ShortURL{Short: strconv.Itoa(s.Counter)}
	s.FullURLKeyMap[fullURL] = sURL
	s.ShortURLKeyMap[sURL.Short] = fURL
	s.Counter++
	return sURL.Short, nil
}
