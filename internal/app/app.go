package app

import (
	hd "github.com/gtgaleevtimur/reduction-url-service/internal/handlers"
	"log"
	"net/http"
)

func Run(addr string) {
	storage := hd.NewServerStore()
	server := &http.Server{Addr: addr, Handler: hd.NewRouter(storage)}
	log.Fatal(server.ListenAndServe())
}
