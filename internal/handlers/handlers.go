package handlers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
	"io"
	"net/http"
)

func NewRouter(controller *ServerStore) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", controller.ReductionURL)
	r.Get("/{id}", controller.GetFullURL)

	return r
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
		buildErrResponse(w, http.StatusMethodNotAllowed, []byte("Need POST requests!"))
		return
	}
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		buildErrResponse(w, http.StatusInternalServerError, []byte("Error via reading request body"))
		return
	}
	inputURL := string(body)
	var result []byte
	shortURL, err := h.Store.Insert(inputURL)
	if err != nil {
		buildErrResponse(w, http.StatusInternalServerError, []byte(err.Error()))
		return
	}
	result = []byte(shortURL)
	w.Header().Set("Content-Type", "text/plain ; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

func (h ServerStore) GetFullURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		buildErrResponse(w, http.StatusMethodNotAllowed, []byte("Need Get requests!"))
		return
	}
	//id := strings.Trim(r.URL.Path, "/")
	id := chi.URLParam(r, "id")
	longURL, err := h.Store.Get(id)
	if err != nil {
		buildErrResponse(w, http.StatusBadRequest, []byte(err.Error()))
		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func buildErrResponse(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write(body)
}
