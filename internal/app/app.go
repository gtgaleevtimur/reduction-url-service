package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"golang.org/x/crypto/acme/autocert"

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
	startServer(conf, handler.NewRouter(storage, conf))
}

// startServer - запускает сервер с настройками из конфигурационного файла.
func startServer(c *config.Config, h chi.Router) {
	if !c.EnableHTTPS {
		server := &http.Server{
			Addr:    c.ServerAddress,
			Handler: h,
		}
		log.Fatal(server.ListenAndServe())
	}
	if c.EnableHTTPS {
		manager := &autocert.Manager{
			Cache:      autocert.DirCache("cache"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(c.ServerAddress),
		}
		server := &http.Server{
			Addr:      ":443",
			Handler:   h,
			TLSConfig: manager.TLSConfig(),
		}
		log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
	}
}
