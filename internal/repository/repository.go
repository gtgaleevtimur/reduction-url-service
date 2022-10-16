package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"sync"
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
	recoveryDisk   string
	sync.Mutex
}

func NewStorage() *Storage {
	s := &Storage{
		Counter:        0,
		FullURLKeyMap:  make(map[string]ShortURL),
		ShortURLKeyMap: make(map[string]FullURL),
	}

	err := s.LoadRecoveryStorage()
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
	if len(s.recoveryDisk) > 0 {
		err = s.AddToRecoveryStorage(&URLItem)
		if err != nil {
			return "", err
		}
	}
	return sURL.Short, nil
}

func (s *Storage) LoadRecoveryStorage() error {
	str, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if !ok {
		return errors.New("env FILE_STORAGE_PATH error")
	}
	file, err := os.OpenFile(str, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	s.recoveryDisk = str
	scanner := bufio.NewScanner(file)
	for {
		if !scanner.Scan() {
			return scanner.Err()
		}
		data := scanner.Bytes()
		URLItem := URL{}
		err = json.Unmarshal(data, &URLItem)
		if err == nil {
			func() {
				s.Lock()
				defer s.Unlock()
				s.Counter++
				s.FullURLKeyMap[URLItem.Full.Full] = URLItem.Short
				s.ShortURLKeyMap[URLItem.Short.Short] = URLItem.Full
			}()
		}
	}
}

func (s *Storage) AddToRecoveryStorage(URLItem *URL) error {
	file, err := os.OpenFile(s.recoveryDisk, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(&URLItem)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	if _, err := writer.Write(data); err != nil {
		return err
	}
	if err := writer.WriteByte('\n'); err != nil {
		return err
	}
	return writer.Flush()
}
