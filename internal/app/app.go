package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	"github.com/gtgaleevtimur/reduction-url-service/internal/handler"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
)

func Run() {
	//Конфигурация приложения через считывание флагов и переменных окружения.
	conf := config.NewConfig(config.WithParseEnv())
	//Инициализация хранилища приложения.
	storage, err := repository.NewDataSource(conf)
	if err != nil {
		log.Fatal(err)
	}
	//Инициализация и запуск сервера.
	server := &http.Server{
		Handler: handler.NewRouter(storage, conf),
		Addr:    conf.ServerAddress,
		BaseContext: func(listener net.Listener) context.Context {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			return ctx
		},
	}
	log.Fatal(server.ListenAndServe())
}
