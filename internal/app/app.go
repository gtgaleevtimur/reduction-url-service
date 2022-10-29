package app

import (
	"log"
	"net/http"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	hd "github.com/gtgaleevtimur/reduction-url-service/internal/handlers"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
)

func Run() {
	//конфигурация приложения через считывание флагов и переменных окружения.
	conf := config.NewConfig(config.WithParseEnv())
	//инициализация хранилища приложения.
	storage := repository.NewStorage(conf)
	//Инициализация и запуск сервера.
	server := &http.Server{
		Handler: hd.NewRouter(storage, conf),
		Addr:    conf.ServerAddress,
	}
	log.Fatal(server.ListenAndServe())

}
