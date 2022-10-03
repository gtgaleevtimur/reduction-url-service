package repository

import (
	"errors"
	"strconv"
	"sync"
)

type URL struct {
	Full  string
	Short string
}

type FullURL struct {
	Full string
}

type ShortURL struct {
	Short string
}

type Storage struct {
	Counter int
	Data    map[int]URL
	sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		Counter: 0,
		Data:    make(map[int]URL),
	}
}

func (s *Storage) GetShortURL(fullURL string) (string, error) {
	s.Lock()
	defer s.Unlock()
	for _, element := range s.Data {
		if element.Full == fullURL {
			return element.Short, nil
		}
	}
	return "", errors.New("wrong URL")
}

func (s *Storage) GetFullURL(shortURL string) (string, error) {
	s.Lock()
	defer s.Unlock()
	for _, element := range s.Data {
		if element.Short == shortURL {
			return element.Full, nil
		}
	}
	return "", errors.New("wrong URL")
}

func (s *Storage) InsertURL(fullURL string) (string, error) {
	if fullURL == "" || fullURL == " " {
		return "", errors.New("ErrNoNilInsert")
	}
	short, err := s.GetShortURL(fullURL)
	if err == nil {
		return short, nil
	}
	s.Lock()
	defer s.Unlock()
	var newURL = URL{Full: fullURL, Short: strconv.Itoa(s.Counter)}
	s.Data[s.Counter] = newURL
	s.Counter++
	return newURL.Short, nil
}
