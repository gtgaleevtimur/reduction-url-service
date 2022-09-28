package app

import (
	hd "github.com/gtgaleevtimur/reduction-url-service/internal/handlers"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
	"log"
	"net/http"
)

var MyStorage *repository.Storage

func Run(addr string) {
	MyStorage = repository.NewStorage()
	mux := http.NewServeMux()
	mux.HandleFunc("/", hd.UrlReduction)
	mux.HandleFunc("/{id}", hd.GetFullUrl)
	server := &http.Server{Addr: addr, Handler: mux}
	log.Fatal(server.ListenAndServe())
}
