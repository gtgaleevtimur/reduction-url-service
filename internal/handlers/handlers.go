package handlers

import (
	"encoding/json"
	"errors"
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
			router.Get("/user/urls", controller.GetAllUserURLs)
			router.Post("/shorten", controller.GetShortURL)
			router.Post("/shorten/batch", controller.PostBatch)

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

//CreateShortURL - обработчик эндпоинта POST / принимает в теле запроса текстовую строку URL для сокращения
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
	statusCode := http.StatusCreated
	//Передаем значения для обработки в хранилище/получаем hash для сокращенного url.
	shortURL, err := h.Storage.MiddlewareInsert(string(textURL), userID.Value)
	if err != nil {
		if errors.Is(err, repository.ErrConflictInsert) {
			statusCode = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	//Создаем сокращенный url.
	exShortURL := h.Conf.ExpShortURL(shortURL)
	//Формируем ответ.
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
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
	userid, err := r.Cookie("shortener")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	statusCode := http.StatusCreated
	var sURL repository.ShortURL
	sURL.Short, err = h.Storage.MiddlewareInsert(full.Full, userid.Value)
	if err != nil {
		if errors.Is(err, repository.ErrConflictInsert) {
			statusCode = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	sURL.Short = h.Conf.ExpShortURL(sURL.Short)
	respBody, err := json.Marshal(sURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
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

func (h ServerHandler) PostBatch(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userid, err := r.Cookie("shortener")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var urls []repository.FullBatch
	if err = json.Unmarshal(body, &urls); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var result []repository.ShortBatch
	for i := range urls {
		short, err := h.Storage.MiddlewareInsert(urls[i].Full, userid.Value)
		if err != nil {
			if errors.Is(err, repository.ErrConflictInsert) {
				result = append(result, repository.ShortBatch{
					Short: h.Conf.ExpShortURL(short),
					CorID: urls[i].CorID,
				})
				continue
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		result = append(result, repository.ShortBatch{
			Short: h.Conf.ExpShortURL(short),
			CorID: urls[i].CorID,
		})
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resultJSON)
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
