package repository

import (
	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
)

// NewDataSource - функция-хэлпер, выбирающая вид хранилища для текущей конфигурации.
func NewDataSource() (result Storager, err error) {
	// Конфигурация приложения через считывание флагов и переменных окружения.
	conf := config.NewConfig(config.WithParseEnv())
	if conf.DatabaseDSN != "" {
		return NewDatabaseDSN(conf)
	} else {
		return NewStorage(conf), nil
	}
}
