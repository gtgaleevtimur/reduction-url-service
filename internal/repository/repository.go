package repository

import (
	"errors"
	"strconv"
	"sync"
)

type Storage struct {
	CountID           int
	IDKeyUrlStorage   map[string]string
	FullUrlKeyStorage map[string]string
	sync.Mutex
}

func New() *Storage {
	return &Storage{
		CountID:           0,
		IDKeyUrlStorage:   make(map[string]string),
		FullUrlKeyStorage: make(map[string]string),
	}

}

func (s *Storage) Insert(fullURL string) (string, error) {
	if fullURL == "" {
		return "", errors.New("ErrEmptyNotAllowed")
	}
	s.Lock()
	defer s.Unlock()
	if value, ok := s.FullUrlKeyStorage[fullURL]; ok {
		return value, nil
	}
	s.IDKeyUrlStorage[strconv.Itoa(s.CountID)] = fullURL
	s.FullUrlKeyStorage[fullURL] = strconv.Itoa(s.CountID)
	s.CountID++
	return s.FullUrlKeyStorage[fullURL], nil
}

func (s *Storage) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("ErrEmptyNotAllowed")
	}
	s.Lock()
	defer s.Unlock()
	if _, ok := s.IDKeyUrlStorage[key]; !ok {
		return "", errors.New("ErrNoKeyStorage")
	}
	return s.IDKeyUrlStorage[key], nil
}
