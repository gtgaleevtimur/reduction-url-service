package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	"github.com/gtgaleevtimur/reduction-url-service/internal/handlers/middlewares"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
)

//NewRouter - функция инициализирующая и настраивающая роутер сервиса.
func NewRouter(s repository.Storager, c *config.Config) chi.Router {
	//инициализация контролера всех хэндлеров приложения.
	controller := newServerHandler(s, c)
	//инициализация роутера chi
	router := chi.NewRouter()
	//запуск поддержки встроенных middleware.
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	//запуск пользовательских middleware.
	router.Use(middleware.Compress(1, `text/plain`, `application/json`))
	router.Use(middleware.AllowContentEncoding(`gzip`))
	router.Use(middlewares.Decompress)
	router.Use(middlewares.CookiesMiddleware)
	//запуск хэндлеров и их паттерны.
	router.Route("/", func(router chi.Router) {
		router.Post("/", controller.CreateShortURL)
		router.Get("/{id}", controller.GetFullURL)
		router.Get("/ping", controller.Ping)

		router.Route("/api", func(router chi.Router) {
			router.Post("/shorten", controller.GetShortURL)
			router.Get("/user/urls", controller.GetAllUserURLs)
		})
	})

	//запуск хэндлеров обработчиков не поддерживаемых методов и маршрутов.
	router.NotFound(NotFound())
	router.MethodNotAllowed(NotAllowed())

	return router
}

//ServerHandler - структура контроллера роутера.
type ServerHandler struct {
	Storage repository.Storager
	Conf    *config.Config
}

//newServerHandler - функция-конструктор контроллера.
func newServerHandler(s repository.Storager, c *config.Config) *ServerHandler {
	return &ServerHandler{Storage: s, Conf: c}
}

//CreateShortURL - обработчик эндпоинта POST / принимает в теле запроса строку URL для сокращения
//и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func (h ServerHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	//Читаем тело и проверяем ошибку.
	textURL, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := r.Cookie("shortener")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Передаем значения для обработки в хранилище/получаем hash для сокращенного url.
	shortURL, err := h.Storage.MiddlewareInsert(string(textURL), userID.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Создаем сокращенный url.
	exShortURL := h.Conf.ExpShortURL(shortURL)
	//Формируем ответ.
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(exShortURL))
}

// GetFullURL -обработчик эндпоинта GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
//и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
func (h ServerHandler) GetFullURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "id")
	if shortURL == "" {
		http.Error(w, "ErrNoEmptyURLParam", http.StatusBadRequest)
		return
	}
	fullURL, err := h.Storage.GetFullURL(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if !strings.HasPrefix(fullURL, config.HTTP) {
		fullURL = config.HTTP + strings.TrimPrefix(fullURL, "//")
	}
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h ServerHandler) GetShortURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var full repository.FullURL
	err = json.Unmarshal(reqBody, &full)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var sURL repository.ShortURL
	sURL.Short, err = h.Storage.GetShortURL(full.Full)
	if err != nil {
		r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		h.insertHelper(w, r)
		return
	}
	sURL.Short = h.Conf.ExpShortURL(sURL.Short)
	respBody, err := json.Marshal(sURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write(respBody)
}

func (h ServerHandler) insertHelper(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := r.Cookie("shortener")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var full repository.FullURL
	err = json.Unmarshal(reqBody, &full)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var sURL repository.ShortURL
	fromInsert, err := h.Storage.MiddlewareInsert(full.Full, userID.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sURL.Short = fromInsert
	sURL.Short = h.Conf.ExpShortURL(sURL.Short)
	respBody, err := json.Marshal(sURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write(respBody)
}

func (h ServerHandler) GetAllUserURLs(w http.ResponseWriter, r *http.Request) {
	userid, err := r.Cookie("shortener")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	urls, err := h.Storage.GetAllUserURLs(userid.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	for i := range urls {
		urls[i].Short = h.Conf.ExpShortURL(urls[i].Short)
	}
	urlsJSON, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(urlsJSON)
}

func (h ServerHandler) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.Storage.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("route does not exist"))
	}
}

func NotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("method does not allowed"))
	}
}
