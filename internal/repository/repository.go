package repository

import (
	"context"
	"errors"
	"io"
	"log"
	"strconv"
	"sync"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
)

type URL struct {
	Full  FullURL
	Short ShortURL
}

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
	FileRecover    *FileRecover
	sync.Mutex
}

func NewStorage(c *config.Config) *Storage {
	s := &Storage{
		Counter:        0,
		FullURLKeyMap:  make(map[string]ShortURL),
		ShortURLKeyMap: make(map[string]FullURL),
	}

	err := s.LoadRecoveryStorage(c.StoragePath)
	if err != nil {
		log.Println(err)
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
	var URLItem = URL{
		Full:  fURL,
		Short: sURL,
	}
	if s.FileRecover != nil {
		err = s.FileRecover.Writer.Write(&URLItem)
		if err != nil {
			return "", err
		}
	}
	return sURL.Short, nil
}

func (s *Storage) LoadRecoveryStorage(str string) error {
	if str == "" {
		return errors.New(" err FILE_STORAGE_PATH is nil ")
	}
	s.Lock()
	defer s.Unlock()
	fileRecover, err := NewFileRecover(str)
	if err != nil {
		return err
	}
	s.FileRecover = fileRecover
	for {
		rURL, err := s.FileRecover.Reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		s.FullURLKeyMap[rURL.Full.Full] = rURL.Short
		s.ShortURLKeyMap[rURL.Short.Short] = rURL.Full
		s.Counter++
	}
	return nil
}
