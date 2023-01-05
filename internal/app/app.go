// Package app аккумулирует все компоненты сервиса и запускает его работу.
package app

import (
	"log"
	"net/http"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	"github.com/gtgaleevtimur/reduction-url-service/internal/handler"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
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
	server := &http.Server{
		Handler: handler.NewRouter(storage, conf),
		Addr:    conf.ServerAddress,
	}
	log.Fatal(server.ListenAndServe())
}
