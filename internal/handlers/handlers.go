package handlers

import (
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
	"io"
	"net/http"
	"strconv"
)

type ServerStore struct {
	Store *repository.Storage
}

func NewServerStore() *ServerStore {
	return &ServerStore{Store: repository.New()}
}

func (h ServerStore) ReductionURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Need POST requests!", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error via reading request body", http.StatusInternalServerError)
		return
	}
	inputURL := string(body)
	var result []byte
	if _, ok := h.Store.FullUrlKeyStorage[inputURL]; !ok {
		h.Store.FullUrlKeyStorage[inputURL] = strconv.Itoa(h.Store.CountID)
		h.Store.IDKeyUrlStorage[strconv.Itoa(h.Store.CountID)] = inputURL
		h.Store.CountID++
	}
	result = []byte(h.Store.FullUrlKeyStorage[inputURL])
	w.Header().Set("Content-Type", "text/plain ; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

func (h ServerStore) GetFullUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Need POST requests!", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Path
	if val, ok := h.Store.IDKeyUrlStorage[id]; !ok {
		http.Error(w, "Dont have URL id in DB", http.StatusBadRequest)
		return
	} else {
		w.Header().Set("Location", val)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
