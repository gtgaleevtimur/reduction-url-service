package app

import (
	hd "github.com/gtgaleevtimur/reduction-url-service/internal/handlers"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
	"log"
	"net/http"
)

func Run(addr string) {
	handler := hd.ServiceHandler{
		Repository: repository.NewStorage(),
	}
	server := &http.Server{Addr: addr, Handler: handler}
	log.Fatal(server.ListenAndServe())
}
