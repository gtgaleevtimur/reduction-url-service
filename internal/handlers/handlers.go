package handlers

import (
	"github.com/gtgaleevtimur/reduction-url-service/internal/app"
	"io"
	"net/http"
	"strconv"
)

func UrlReduction(w http.ResponseWriter, r *http.Request) {
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
	if _, ok := app.MyStorage.FullUrlKeyStorage[inputURL]; !ok {
		app.MyStorage.FullUrlKeyStorage[inputURL] = strconv.Itoa(app.MyStorage.CountID)
		app.MyStorage.IDKeyUrlStorage[strconv.Itoa(app.MyStorage.CountID)] = inputURL
		app.MyStorage.CountID++
	}
	result = []byte(app.MyStorage.FullUrlKeyStorage[inputURL])
	w.Header().Set("Content-Type", "text/plain ; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

func GetFullUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Need POST requests!", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Path
	if val, ok := app.MyStorage.IDKeyUrlStorage[id]; !ok {
		http.Error(w, "Dont have URL id in DB", http.StatusBadRequest)
		return
	} else {
		w.Header().Set("Location", val)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
