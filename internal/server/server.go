package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"golang.org/x/crypto/acme/autocert"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
)

//RunServer - запускает сервер с настройками из конфигурационного файла.
func RunServer(c *config.Config, h chi.Router) {
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
