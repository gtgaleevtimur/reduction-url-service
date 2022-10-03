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

func New() *Storage {
	return &Storage{
		Counter: 0,
		Data:    make(map[int]URL),
	}
}

func (s *Storage) GetShortURL(fullURL string) (*ShortURL, error) {
	s.Lock()
	defer s.Unlock()
	for _, element := range s.Data {
		if element.Full == fullURL {
			return &ShortURL{Short: element.Short}, nil
		}
	}
	return nil, errors.New("wrong URL")
}

func (s *Storage) GetFullURL(shortURL string) (*FullURL, error) {
	s.Lock()
	defer s.Unlock()
	for _, element := range s.Data {
		if element.Short == shortURL {
			return &FullURL{Full: element.Full}, nil
		}
	}
	return nil, errors.New("wrong URL")
}

func (s *Storage) InsertURL(fullURL string) (string, error) {
	short, err := s.GetShortURL(fullURL)
	if err == nil {
		return short.Short, nil
	}
	s.Lock()
	defer s.Unlock()
	var newURL = URL{Full: fullURL, Short: "/" + strconv.Itoa(s.Counter)}
	s.Data[s.Counter] = newURL
	s.Counter++
	return newURL.Short, nil
}
