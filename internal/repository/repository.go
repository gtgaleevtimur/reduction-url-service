package repository

import (
	"crypto/sha1"
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

//NewStorage - функция-конструктор in-memory хранилища,возвращает интерфейс.
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
func (s *Storage) MiddlewareInsert(fullURL string, userID string) (string, error) {
	//Генерируем hash.
	hasher := sha1.Sum([]byte(fullURL + userID))
	hash := hex.EncodeToString(hasher[:len(hasher)/5])
	//Проверяем есть ли в хранилище такой url.
	okHash, err := s.GetShortURL(fullURL)
	//Если нет,то вставляем новые данные.
	if err != nil {
		err = s.InsertURL(fullURL, userID, hash)
		if err != nil {
			return "", err
		}
		//Возвращаем сгенерированный hash.
		return hash, nil
	}
	//Если есть , возвращаем hash и ошибку.
	return okHash, ErrConflictInsert
}

//GetShortURL - метод, возвращающий hash сокращенного url.
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

//GetFullURL - метод, возвращающий original_url по его hash.
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
	//Проверяем полученные данные.
	if fullURL == "" || fullURL == " " || userid == "" || userid == " " || hash == "" || hash == " " {
		return errors.New("ErrNoEmptyInsert")
	}
	//Блокируем хранилище на время операции.
	s.Lock()
	defer s.Unlock()
	//Записываем данные в хранилище.
	s.Data[hash] = URL{
		UserID: userid,
		FURL:   fullURL,
	}
	//Если FILE_STORAGE_PATH выставлен,то записывает данные в резервное хранилище..
	if s.FileRecover != nil {
		//Готовим структуру для резервного хранилища.
		URLItem := NodeURL{
			Hash:   hash,
			FURL:   fullURL,
			UserID: userid,
		}
		//Записываем.
		err := s.FileRecover.Writer.Write(&URLItem)
		if err != nil {
			return err
		}
	}
	return nil
}

//LoadRecoveryStorage - метод , восстанавливающий данные из резервного хранилища при инициализации in-memory.
func (s *Storage) LoadRecoveryStorage(str string) error {
	//Выполняем проверку текущей конфигурации.
	if str == "" {
		return errors.New("err FILE_STORAGE_PATH is nil ")
	}
	//Блокируем хранилище на время выполнения операции.
	s.Lock()
	defer s.Unlock()
	//Создаем FileRecover.
	fileRecover, err := NewFileRecover(str)
	if err != nil {
		return err
	}
	s.FileRecover = fileRecover
	for {
		//Читаем построчно из резервного хранилища данные.
		node, err := s.FileRecover.Reader.Read()
		//Проверка ошибки.
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		//Вставляем считанные данные.
		s.Data[node.Hash] = URL{
			UserID: node.UserID,
			FURL:   node.FURL,
		}
	}
	return nil
}

//GetAllUserURLs - метод возвращающий массив со всеми original_url+hash сохраненными пользователем.
func (s *Storage) GetAllUserURLs(userid string) ([]SlicedURL, error) {
	//Блокируем хранилище на время выполнения операции.
	s.Lock()
	defer s.Unlock()
	//Инициализируем результирующий массив.
	result := make([]SlicedURL, 0)
	//Итерируемся по хранилищу
	for hash, url := range s.Data {
		//Если нашли совпадение userID,то добавляем в массив данные.
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

//Ping - метод заглушка для in-memory.
func (s *Storage) Ping() error {
	return nil
}
