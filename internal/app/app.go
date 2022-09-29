package app

import (
	hd "github.com/gtgaleevtimur/reduction-url-service/internal/handlers"
	"log"
	"net/http"
)

func Run(addr string) {
	mux := http.NewServeMux()
	store := hd.NewServerStore()
	mux.HandleFunc("/", store.Root)
	server := &http.Server{Addr: addr, Handler: mux}
	log.Fatal(server.ListenAndServe())
}
