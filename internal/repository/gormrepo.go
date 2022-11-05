package repository

import (
	"crypto/md5"
	"encoding/hex"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
)

type GormURL struct {
	Hash   string
	UserID string
	FURL   string
}

type GormDatabase struct {
	DB *gorm.DB
}

func NewGormStorage(conf *config.Config) (Storager, error) {
	db, err := gorm.Open(postgres.Open(conf.DatabaseDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&GormURL{})
	if err != nil {
		return nil, err
	}
	s := &GormDatabase{
		DB: db,
	}
	return s, nil
}

func (g GormDatabase) GetShortURL(fullURL string) (string, error) {
	var gURL GormURL
	g.DB.Where("furl = ?", fullURL).First(&gURL)
	return gURL.Hash, nil
}

func (g GormDatabase) GetFullURL(shortURL string) (string, error) {
	var gURL GormURL
	g.DB.Where("hash = ?", shortURL).First(&gURL)
	return gURL.FURL, nil
}

func (g GormDatabase) InsertURL(fullURL string, userid string, hash string) error {
	gURL := GormURL{
		FURL:   fullURL,
		UserID: userid,
		Hash:   hash,
	}
	g.DB.Create(&gURL)
	return nil
}

func (g GormDatabase) MiddlewareInsert(fullURL string, userID string) (string, error) {
	hasher := md5.Sum([]byte(fullURL + userID))
	hash := hex.EncodeToString(hasher[:len(hasher)/5])
	okHash, _ := g.GetShortURL(fullURL)
	//Если нет,то вставляем новые данные.
	if okHash != "" {
		_ = g.InsertURL(fullURL, userID, hash)
		//Возвращаем сгенерированный hash.
		return hash, nil
	}
	//Если есть , возвращаем hash и ошибку.
	return okHash, ErrConflictInsert
}

func (g GormDatabase) GetAllUserURLs(userid string) ([]SlicedURL, error) {
	var gURLs []GormURL
	g.DB.Where("userid = ?", userid).Find(&gURLs)
	var result []SlicedURL
	for _, value := range gURLs {
		res := SlicedURL{
			Full:  value.FURL,
			Short: value.Hash,
		}
		result = append(result, res)
	}
	return result, nil
}

func (g GormDatabase) Ping() error {
	return nil
}
