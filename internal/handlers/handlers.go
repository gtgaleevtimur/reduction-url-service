package handlers

import (
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
	"io"
	"net/http"
	"strconv"
)

type ServiceHandler struct {
	Repository *repository.Storage
}

func (h ServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error via reading request body", http.StatusInternalServerError)
			return
		}
		inputURL := string(body)
		var result []byte
		if _, ok := h.Repository.FullUrlKeyStorage[inputURL]; !ok {
			h.Repository.FullUrlKeyStorage[inputURL] = strconv.Itoa(h.Repository.CountID)
			h.Repository.IDKeyUrlStorage[strconv.Itoa(h.Repository.CountID)] = inputURL
			h.Repository.CountID++
		}
		result = []byte(h.Repository.FullUrlKeyStorage[inputURL])
		w.Header().Set("Content-Type", "text/plain ; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write(result)
	}
	if r.Method == http.MethodGet {
		id := r.URL.Path
		if val, ok := h.Repository.IDKeyUrlStorage[id]; !ok {
			http.Error(w, "Dont have URL id in DB", http.StatusBadRequest)
			return
		} else {
			w.Header().Set("Location", val)
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
	}
	http.Error(w, "Need POST requests!", http.StatusMethodNotAllowed)
}
