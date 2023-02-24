package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gtgaleevtimur/reduction-url-service/proto"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	"github.com/gtgaleevtimur/reduction-url-service/internal/grpcserv"
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
	startServer(conf, storage)
}

// startServer - запускает сервер с настройками из конфигурационного файла.
func startServer(conf *config.Config, storage repository.Storager) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcserv.MyUnaryInterceptor))

	if !conf.EnableHTTPS {
		if conf.EnableGRPC {
			go startGRPC(storage, conf, grpcServer, cancel)
		}
		server := &http.Server{
			Addr:    conf.ServerAddress,
			Handler: handler.NewRouter(storage, conf),
		}

		go gracefulShutdown(ctx, server, grpcServer)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			cancel()
			log.Fatal("failed to run server")
		}
	}
	if conf.EnableHTTPS {
		if conf.EnableGRPC {
			go startGRPC(storage, conf, grpcServer, cancel)
		}
		manager := &autocert.Manager{
			Cache:      autocert.DirCache("cache"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(conf.ServerAddress),
		}
		server := &http.Server{
			Addr:      ":443",
			Handler:   handler.NewRouter(storage, conf),
			TLSConfig: manager.TLSConfig(),
		}

		go gracefulShutdown(ctx, server, grpcServer)

		err := server.ListenAndServeTLS("server.crt", "server.key")
		if err != nil && err != http.ErrServerClosed {
			cancel()
			log.Fatal("failed to run server")
		}
	}
}

// gracefulShutdown - GracefulShutdown по сигналу syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT.
func gracefulShutdown(ctx context.Context, server *http.Server, grpcServer *grpc.Server) {
	<-ctx.Done()
	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*20)
	defer shutdownCtxCancel()
	grpcServer.Stop()
	log.Println("gRPC server shutdown")
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

// startGRPC - запуск grpc сервера.
func startGRPC(storage repository.Storager, conf *config.Config, grpcServer *grpc.Server, cancel context.CancelFunc) {
	listen, err := net.Listen("tcp", ":0")
	if err != nil {
		cancel()
		log.Fatal(err.Error())
	}
	proto.RegisterShortenerServer(grpcServer, grpcserv.New(storage, conf))
	log.Println("gRPC server start at:", listen.Addr().String())
	if err = grpcServer.Serve(listen); err != nil {
		cancel()
		log.Fatal(err.Error())
	}
}
