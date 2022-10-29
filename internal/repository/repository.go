package repository

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"sync"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
)

//Storage - структура in-memory хранилища.
type Storage struct {
	Data        map[string]URL
	FileRecover *FileRecover
	sync.Mutex
}

//NewStorage - функция-конструктор in-memory хранилища.
func NewStorage(c *config.Config) Storager {
	s := &Storage{
		Data: make(map[string]URL),
	}

	//Проверяем задан ли FILE_STORAGE_PATH,если да,то восстанавливаем данные оттуда.
	err := s.LoadRecoveryStorage(c.StoragePath)
	if err != nil {
		log.Println(err)
	}

	return s
}

//MiddlewareInsert - метод-помощник, генерирует hash для ключа,передает hash+url+userid хранилищу,возвращает сокращенный url
func (s *Storage) MiddlewareInsert(fURL string, userID string) (string, error) {
	//Генерируем hash.
	hasher := md5.Sum([]byte(fURL + userID))
	hash := hex.EncodeToString(hasher[:len(hasher)/5])
	//Проверяем есть ли в хранилище такой url.
	okHash, err := s.GetShortURL(fURL)
	//Если нет,то вставляем новые данные.
	if err != nil {
		err = s.InsertURL(fURL, userID, hash)
		if err != nil {
			return "", err
		}
		//Возвращаем сгенерированный hash.
		return hash, nil
	}
	//Если есть , возвращаем hash.
	return okHash, nil
}

//GetShortURL - метод-помощник,возвращает hash url если полный url есть в хранилище.
func (s *Storage) GetShortURL(fullURL string) (string, error) {
	s.Lock()
	defer s.Unlock()
	for hash, value := range s.Data {
		if value.FURL == fullURL {
			return hash, nil
		}
	}
	return "", errors.New("ErrNotFoundURL")
}

//GetFullURL - возвращает полный url по hash сокращенного.
func (s *Storage) GetFullURL(shortURL string) (string, error) {
	s.Lock()
	defer s.Unlock()
	if val, ok := s.Data[shortURL]; ok {
		return val.FURL, nil
	}
	return "", errors.New("ErrNotFoundURL")
}

//InsertURL - метод,заполняющий хранилище данными(полный url, id пользователя).
func (s *Storage) InsertURL(fullURL string, userid string, hash string) error {
	if fullURL == "" || fullURL == " " || userid == "" || userid == " " || hash == "" || hash == " " {
		return errors.New("ErrNoEmptyInsert")
	}
	s.Lock()
	defer s.Unlock()
	s.Data[hash] = URL{
		UserID: userid,
		FURL:   fullURL,
	}
	var URLItem = NodeURL{
		Hash:   hash,
		FURL:   fullURL,
		UserID: userid,
	}
	//Если FILE_STORAGE_PATH выставлен,то записывает данные туда.
	if s.FileRecover != nil {
		err := s.FileRecover.Writer.Write(&URLItem)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) LoadRecoveryStorage(str string) error {
	if str == "" {
		return errors.New("err FILE_STORAGE_PATH is nil ")
	}
	s.Lock()
	defer s.Unlock()
	fileRecover, err := NewFileRecover(str)
	if err != nil {
		return err
	}
	s.FileRecover = fileRecover
	for {
		node, err := s.FileRecover.Reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		s.Data[node.Hash] = URL{
			UserID: node.UserID,
			FURL:   node.FURL,
		}
	}
	return nil
}

func (s *Storage) GetAllUserURLs(userid string) ([]SlicedURL, error) {
	s.Lock()
	defer s.Unlock()

	result := make([]SlicedURL, 0)

	for hash, url := range s.Data {
		if url.UserID == userid {
			result = append(result, SlicedURL{
				Short: hash,
				Full:  url.FURL,
			})
		}
	}
	if len(result) == 0 {
		return nil, errors.New("ErrNotExistUserURLs")
	} else {
		return result, nil
	}
}
