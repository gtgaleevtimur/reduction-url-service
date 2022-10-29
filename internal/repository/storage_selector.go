package repository

import (
	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
)

func NewDataSource(conf *config.Config) (result Storager, err error) {
	if conf.DatabaseDSN != "" {
		result, err = NewDatabaseDSN(conf)
		return NewDatabaseDSN(conf)
	} else {
		NewStorage(conf)
		return NewStorage(conf), nil
	}
}
