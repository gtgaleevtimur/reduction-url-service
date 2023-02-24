package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if !c.EnableHTTPS {
		server := &http.Server{
			Addr:    c.ServerAddress,
			Handler: h,
		}

		go gracefulShutdown(server, sig)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to run server")
		}
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

		go gracefulShutdown(server, sig)

		err := server.ListenAndServeTLS("server.crt", "server.key")
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to run server")
		}
	}
}

// gracefulShutdown - GracefulShutdown по сигналу syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT.
func gracefulShutdown(server *http.Server, sig chan os.Signal) {
	<-sig
	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*20)
	defer shutdownCtxCancel()
	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			log.Fatal("graceful shutdown timed out and forcing exit.")
		}
	}()
	err := server.Shutdown(context.Background())
	if err != nil {
		log.Fatal("server shutdown error")
	}
}
