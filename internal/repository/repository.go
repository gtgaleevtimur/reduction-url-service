package repository

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"sync"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
)

// Storage - структура in-memory хранилища.
type Storage struct {
	Data        map[string]URL
	FileRecover *FileRecover
	sync.Mutex
}

// NewStorage - функция-конструктор in-memory хранилища,возвращает интерфейс.
func NewStorage(c *config.Config) Storager {
	s := &Storage{
		Data: make(map[string]URL),
	}

	// Проверяем задан ли FILE_STORAGE_PATH, если да, то восстанавливаем данные оттуда.
	err := s.LoadRecoveryStorage(c.StoragePath)
	if err != nil {
		if !errors.Is(err, ErrFileStoragePathNil) {
			log.Println(err)
		}
	}

	return s
}

// InsertURL - метод ,который генерирует hash для ключа,передает hash+url+userid хранилищу,возвращает сокращенный url
func (s *Storage) InsertURL(ctx context.Context, fullURL string, userID string) (string, error) {
	// Генерируем hash.
	hasher := md5.Sum([]byte(fullURL + userID))
	hash := hex.EncodeToString(hasher[:len(hasher)/5])
	// Проверяем есть ли в хранилище такой url.
	okHash, err := s.GetShortURL(ctx, fullURL)
	// Если нет, то вставляем новые данные.
	if err != nil {
		err = s.saveData(ctx, fullURL, userID, hash)
		if err != nil {
			return "", err
		}
		// Возвращаем сгенерированный hash.
		return hash, nil
	}
	// Если есть, возвращаем hash и ошибку.
	return okHash, ErrConflictInsert
}

// GetShortURL - метод, возвращающий hash сокращенного url.
func (s *Storage) GetShortURL(_ context.Context, fullURL string) (string, error) {
	s.Lock()
	defer s.Unlock()
	for hash, value := range s.Data {
		if value.FURL == fullURL && !value.Delete {
			return hash, nil
		}
	}
	return "", ErrNotFoundURL
}

// GetFullURL - метод, возвращающий original_url по его hash.
func (s *Storage) GetFullURL(_ context.Context, shortURL string) (string, error) {
	s.Lock()
	defer s.Unlock()
	// Проверяем наличие URL в БД.
	if val, ok := s.Data[shortURL]; ok {
		// Если URL удален возвращаем соответствующую ошибку.
		if val.Delete {
			return "", ErrDeletedURL
		}
		// Возвращаем URL если все ок.
		return val.FURL, nil
	}
	// Если URL отсутствует в БД возвращаем соответствующую ошибку.
	return "", ErrNotFoundURL
}

// InsertURL - метод,заполняющий хранилище данными(полный url, id пользователя, hash).
func (s *Storage) saveData(_ context.Context, fullURL string, userid string, hash string) error {
	// Проверяем полученные данные.
	if fullURL == "" || fullURL == " " || userid == "" || userid == " " || hash == "" || hash == " " {
		return errors.New("ErrNoEmptyInsert")
	}
	// Блокируем хранилище на время операции.
	s.Lock()
	defer s.Unlock()
	// Записываем данные в хранилище.
	s.Data[hash] = URL{
		UserID: userid,
		FURL:   fullURL,
		Delete: false,
	}
	// Если FILE_STORAGE_PATH выставлен, нто записывает данные в резервное хранилище..
	if s.FileRecover != nil {
		// Готовим структуру для резервного хранилища.
		URLItem := NodeURL{
			Hash:   hash,
			FURL:   fullURL,
			UserID: userid,
			Delete: false,
		}
		// Записываем.
		err := s.FileRecover.Writer.Write(&URLItem)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadRecoveryStorage - метод, восстанавливающий данные из резервного хранилища при инициализации in-memory.
func (s *Storage) LoadRecoveryStorage(str string) error {
	// Выполняем проверку текущей конфигурации.
	if str == "" {
		return ErrFileStoragePathNil
	}
	// Блокируем хранилище на время выполнения операции.
	s.Lock()
	defer s.Unlock()
	// Создаем FileRecover.
	fileRecover, err := NewFileRecover(str)
	if err != nil {
		return err
	}
	s.FileRecover = fileRecover
	for {
		// Читаем построчно из резервного хранилища данные.
		node, err := s.FileRecover.Reader.Read()
		// Проверка ошибки.
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		// Вставляем считанные данные.
		s.Data[node.Hash] = URL{
			UserID: node.UserID,
			FURL:   node.FURL,
			Delete: node.Delete,
		}
	}
	return nil
}

// GetAllUserURLs - метод возвращающий массив со всеми original_url+hash сохраненными пользователем.
func (s *Storage) GetAllUserURLs(_ context.Context, userid string) ([]SlicedURL, error) {
	// Блокируем хранилище на время выполнения операции.
	s.Lock()
	defer s.Unlock()
	// Инициализируем результирующий массив.
	result := make([]SlicedURL, 0)
	// Итерируемся по хранилищу
	for hash, url := range s.Data {
		// Если нашли совпадение userID, то добавляем в массив данные.
		if url.UserID == userid && !url.Delete {
			result = append(result, SlicedURL{
				Short: hash,
				Full:  url.FURL,
			})
		}
	}
	// Если записи не найдены возвращаем ошибку
	if len(result) == 0 {
		return nil, errors.New("ErrNotExistUserURLs")
	}
	// Иначе возвращаем массив
	return result, nil
}

// Delete - метод, который данные помечает как удаленные по их hash(идентификатор).
func (s *Storage) Delete(_ context.Context, hashes []string, userID string) error {
	// Блокируем хранилище на время выполнения операции.
	s.Lock()
	defer s.Unlock()
	// Проверяем что userID URL в базе данных с таким hash соответствует userID, сделавшему запрос
	for _, hash := range hashes {
		if s.Data[hash].UserID == userID {
			// Применяем изменения.
			s.Data[hash] = URL{
				UserID: userID,
				FURL:   s.Data[hash].FURL,
				Delete: true,
			}
			// Если задан файл для резервного хранения, то пишем так же туда.
			if s.FileRecover != nil {
				URLItem := NodeURL{
					Hash:   hash,
					FURL:   s.Data[hash].FURL,
					UserID: userID,
					Delete: true,
				}
				// Записываем.
				err := s.FileRecover.Writer.Write(&URLItem)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Ping - метод заглушка для in-memory.
func (s *Storage) Ping(_ context.Context) error {
	return nil
}
