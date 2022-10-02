package handlers

import (
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
	"io"
	"net/http"
	"strings"
)

type MyHandlers interface {
	ReductionURL(w http.ResponseWriter, r *http.Request)
	GetFullURL(w http.ResponseWriter, r *http.Request)
	Root(w http.ResponseWriter, r *http.Request)
}

type ServerStore struct {
	Store *repository.Storage
}

func NewServerStore() *ServerStore {
	return &ServerStore{Store: repository.New()}
}

func (h ServerStore) Root(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if len(path) == 1 {
		h.ReductionURL(w, r)
	} else {
		h.GetFullURL(w, r)
	}
}

func (h ServerStore) ReductionURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Need POST requests!", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error via reading request body", http.StatusInternalServerError)
		return
	}
	inputURL := string(body)
	var result []byte
	shortURL, err := h.Store.Insert(inputURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result = []byte(shortURL)
	w.Header().Set("Content-Type", "text/plain ; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

func (h ServerStore) GetFullURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Need Get requests!", http.StatusMethodNotAllowed)
		return
	}
	id := strings.Trim(r.URL.Path, "/")
	longURL, err := h.Store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
