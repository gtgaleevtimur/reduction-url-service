package repository

import (
	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
)

//NewDataSource - функция-хэлпер,выбирающая вид хранилища для текущей конфигурации.
func NewDataSource(conf *config.Config) (result Storager, err error) {
	if conf.DatabaseDSN != "" {
		return NewDatabaseDSN(conf)
	} else {
		return NewStorage(conf), nil
	}
}
