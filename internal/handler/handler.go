package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
)

// NewRouter - функция инициализирующая и настраивающая роутер сервиса.
func NewRouter(s repository.Storager, c *config.Config) chi.Router {
	// Инициализация контролера всех хэндлеров приложения.
	controller := newServerHandler(s, c)
	// Инициализация роутера chi.
	router := chi.NewRouter()
	// Запуск поддержки встроенных middleware.
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	// Запуск пользовательских middleware.
	router.Use(middleware.Compress(1, `text/plain`, `application/json`))
	router.Use(middleware.AllowContentEncoding(`gzip`))
	router.Use(Decompress)
	router.Use(CookiesMiddleware)
	// Запуск хэндлеров и их паттерны.
	router.Route("/", func(router chi.Router) {
		router.Post("/", controller.ShortURLTextBy)
		router.Get("/{hash}", controller.FullURLHashBy)
		router.Get("/ping", controller.Ping)

		router.Route("/api", func(router chi.Router) {
			router.Get("/internal/stats", controller.GetStats)
			router.Delete("/user/urls", controller.DeleteBatch)
			router.Get("/user/urls", controller.GetAllUserURLs)
			router.Post("/shorten", controller.ShortURLJSONBy)
			router.Post("/shorten/batch", controller.PostBatch)
		})
	})
	// Запуск хэндлеров обработчиков не поддерживаемых методов и маршрутов.
	router.NotFound(NotFound())
	router.MethodNotAllowed(NotAllowed())

	return router
}

// ServerHandler - структура контроллера роутера.
type ServerHandler struct {
	Storage repository.Storager
	Conf    *config.Config
}

// newServerHandler - конструктор контроллера.
func newServerHandler(s repository.Storager, c *config.Config) *ServerHandler {
	return &ServerHandler{Storage: s, Conf: c}
}

// GetStats - обработчик эндпоинта GET /api/internal/stats , проверяет реальный IP возвращает статистику по сокращенным
// URL и пользователям в системе.
func (h ServerHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	if h.Conf.TrustedSubnet != "" {
		_, ipNet, err := net.ParseCIDR(h.Conf.TrustedSubnet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		ip, err := GetIP(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if !ipNet.Contains(ip) {
			http.Error(w, repository.ErrCIDRContain.Error(), http.StatusForbidden)
			return
		}
	}
	// Инициализируем контекст
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	urls, err := h.Storage.GetCountURL(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	users, err := h.Storage.GetCountUsers(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	temp := &repository.StatStruct{
		Users: users,
		Urls:  urls,
	}
	response, err := json.Marshal(temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// ShortURLTextBy - обработчик эндпоинта POST /, принимает в теле запроса текстовую строку URL для сокращения.
// Возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func (h ServerHandler) ShortURLTextBy(w http.ResponseWriter, r *http.Request) {
	// Инициализируем контекст
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	// Читаем тело и проверяем ошибку.
	textURL, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Считываем cookie пользователя.
	userID, err := r.Cookie("shortener")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	statusCode := http.StatusCreated
	// Передаем полученные значения для обработки в хранилище/получаем hash сокращенного url.
	hash, err := h.Storage.InsertURL(ctx, string(textURL), userID.Value)
	if err != nil {
		// Проверяем ошибку на соответсвие ситуации, когда вносимый URL уже в базе данных.
		if errors.Is(err, repository.ErrConflictInsert) {
			statusCode = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	// Создаем сокращенный url для возврата пользователю.
	exShortURL := h.Conf.ExpShortURL(hash)
	// Формируем ответ.
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write([]byte(exShortURL))
}

// FullURLHashBy -обработчик эндпоинта GET /{id} ,принимает в качестве URL-параметра идентификатор сокращённого URL.
// Возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
func (h ServerHandler) FullURLHashBy(w http.ResponseWriter, r *http.Request) {
	// Инициализируем контекст.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	// Считываем hash сокращенного URL из параметров запроса.
	shortURL := chi.URLParam(r, "hash")
	if shortURL == "" {
		http.Error(w, "ErrNoEmptyURLParam", http.StatusBadRequest)
		return
	}
	// Запрашиваем оригинальный URL из базы данных.
	fullURL, err := h.Storage.GetFullURL(ctx, shortURL)
	if err != nil {
		if errors.Is(err, repository.ErrDeletedURL) {
			w.WriteHeader(http.StatusGone)
			return
		}
		http.Error(w, "NotExistURL", http.StatusNotFound)
		return
	}
	if !strings.HasPrefix(fullURL, config.HTTP) {
		fullURL = config.HTTP + strings.TrimPrefix(fullURL, "//")
	}
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// ShortURLJSONBy - обработчик эндпоинта POST /api/shorten,принимает в теле запроса json с оригинальным URL.
// Возвращает JSON с сокращенным URL.
func (h ServerHandler) ShortURLJSONBy(w http.ResponseWriter, r *http.Request) {
	// Инициализируем контекст.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	// Читаем тело запроса.
	reqBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Десерриализуем тело запроса в структуру оригинального URL.
	var full repository.FullURL
	err = json.Unmarshal(reqBody, &full)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Считываем cookie пользователя.
	userid, err := r.Cookie("shortener")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	statusCode := http.StatusCreated
	// Готовим структуру с сокращенным URL.
	var sURL repository.ShortURL
	// Передаем данные для сохранения/проверки на сокхранение URL методу базы данных.
	// Возвращает хэш сохраненного URL.
	sURL.Short, err = h.Storage.InsertURL(ctx, full.Full, userid.Value)
	if err != nil {
		// Проверяем ошибку на соответсвие ситуации, когда вносимый URL уже в базе данных.
		if errors.Is(err, repository.ErrConflictInsert) {
			statusCode = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	// Создаем сокращенный url для возврата пользователю.
	sURL.Short = h.Conf.ExpShortURL(sURL.Short)
	// Сериализуем готовую структуру в JSON.
	respBody, err := json.Marshal(sURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Формируем ответ.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(respBody)
}

// GetAllUserURLs - обработчик эндпоинта GET /api/user/urls, считывая userid из cookie возвращает все URL пользователя.
func (h ServerHandler) GetAllUserURLs(w http.ResponseWriter, r *http.Request) {
	// Инициализируем контекст.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	// Считываем cookie пользователя.
	userid, err := r.Cookie("shortener")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Готовим массив с hash сохраненных URL пользвателя.
	urls, err := h.Storage.GetAllUserURLs(ctx, userid.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}
	// Формируем поля с сокращенными URL.
	for i := range urls {
		urls[i].Short = h.Conf.ExpShortURL(urls[i].Short)
	}
	// Серриализуем полученный массив со структурами в JSON.
	urlsJSON, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Формируем ответ.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(urlsJSON)
}

// Ping - обработчик эндпоинта GET /ping , отражает доступность базы данных.
func (h ServerHandler) Ping(w http.ResponseWriter, r *http.Request) {
	// Инициализируем контекст.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	err := h.Storage.Ping(ctx)
	statusCode := http.StatusOK
	if err != nil {
		statusCode = http.StatusInternalServerError
	}
	w.WriteHeader(statusCode)
}

// PostBatch - обработчик эндпоинта POST /api/shorten/batch , принимает в теле запроса массив с JSON.
// Возвращет массив с JSON .
func (h ServerHandler) PostBatch(w http.ResponseWriter, r *http.Request) {
	// Инициализируем контекст.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	// Читаем тело запроса.
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Получаем cookie пользователя.
	userid, err := r.Cookie("shortener")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Готовим массив со структурами и десерриализуем в него тело запроса.
	var urls []repository.FullBatch
	if err = json.Unmarshal(body, &urls); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Готовим массив со структурами для ответа.
	var result []repository.ShortBatch
	// Итерируемся по массиву с полученными данными и сохраняем в базу данных.
	for i := range urls {
		short, err := h.Storage.InsertURL(ctx, urls[i].Full, userid.Value)
		if err != nil {
			if errors.Is(err, repository.ErrConflictInsert) {
				// Заполняем массив с ответом в случае соответсвия ошибки.
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
		// Заполняем массив с ответом.
		result = append(result, repository.ShortBatch{
			Short: h.Conf.ExpShortURL(short),
			CorID: urls[i].CorID,
		})
	}
	// Серриализуем массив с ответом в JSON.
	resultJSON, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Формируем ответ.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resultJSON)
}

// DeleteBatch - обработчик эндпоинта DELETE /api/user/urls , принимает в теле запроса JSON.
// Запускает асинхронный процесс удаления этих URL.
func (h ServerHandler) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	// Инициализируем контекст.
	ctx := context.Background()
	defer r.Body.Close()
	// Читаем тело запроса.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Получаем cookie пользователя.
	userid, err := r.Cookie("shortener")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Создаем массив для разбора тела запроса.
	var hashes []string
	// Так как ожидаем в теле запроса массив строк [ "a", "b", "c", "d", ...].
	// Парсим запрос и записываем результат в массив для разбора.
	err = json.Unmarshal(body, &hashes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// В отдельной горутине запускаем процесс удаления.
	// Передаем горутине список и cookie.
	go h.Storage.Delete(ctx, hashes, userid.Value)
	// Пишем ответ.
	w.WriteHeader(http.StatusAccepted)
}

// NotFound - обработчик неподдерживаемых маршрутов.
func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("route does not exist"))
	}
}

// NotAllowed - обработчик неподдерживаемых методов.
func NotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("method does not allowed"))
	}
}
