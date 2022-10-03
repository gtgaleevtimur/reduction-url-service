package repository

import (
	"errors"
	"strconv"
	"sync"
)

type Storage struct {
	CountID           int
	IDKeyURLStorage   map[string]string
	FullURLKeyStorage map[string]string
	sync.Mutex
}

func New() *Storage {
	return &Storage{
		CountID:           0,
		IDKeyURLStorage:   make(map[string]string),
		FullURLKeyStorage: make(map[string]string),
	}
}

func (s *Storage) Insert(fullURL string) (string, error) {
	if fullURL == "" {
		return "", errors.New("ErrEmptyNotAllowed")
	}
	s.Lock()
	defer s.Unlock()
	if value, ok := s.FullURLKeyStorage[fullURL]; ok {
		return value, nil
	}
	s.IDKeyURLStorage[strconv.Itoa(s.CountID)] = fullURL
	s.FullURLKeyStorage[fullURL] = strconv.Itoa(s.CountID)
	s.CountID++
	return s.FullURLKeyStorage[fullURL], nil
}

func (s *Storage) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("ErrEmptyNotAllowed")
	}
	s.Lock()
	defer s.Unlock()
	if _, ok := s.IDKeyURLStorage[key]; !ok {
		return "", errors.New("ErrNoKeyStorage")
	}
	return s.IDKeyURLStorage[key], nil
}
