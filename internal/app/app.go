// Package app аккумулирует все компоненты сервиса и запускает его работу.
package app

import (
	"log"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	"github.com/gtgaleevtimur/reduction-url-service/internal/handler"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
	"github.com/gtgaleevtimur/reduction-url-service/internal/server"
)

// Run - функция собирающая все компоненты сервиса воедино.
func Run() {
	// Конфигурация приложения через считывание флагов и переменных окружения.
	conf := config.NewConfig(config.WithParseEnv())
	// Инициализация хранилища приложения.
	storage, err := repository.NewDataSource(conf)
	if err != nil {
		log.Fatal(err)
	}
	// Инициализация и запуск сервера.
	server.RunServer(conf, handler.NewRouter(storage, conf))
}
