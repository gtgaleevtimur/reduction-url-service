package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
	"io/ioutil"
	"net/http"
	"strings"
)

func NewRouter(s *repository.Storage) *gin.Engine {
	controller := NewServerHandler(s)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	api := router.Group("/api")
	api.POST("/shorten", controller.GetShortURL)
	router.POST("/", controller.CreateShortURL)
	router.GET("/:id", controller.GetFullURL)
	router.NoRoute(controller.ResponseBadRequest)
	return router
}

type ServerHandler struct {
	Storage *repository.Storage
}

func NewServerHandler(s *repository.Storage) *ServerHandler {
	return &ServerHandler{Storage: s}
}

func (h ServerHandler) CreateShortURL(c *gin.Context) {
	fullURL, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	shortURL, err := h.Storage.InsertURL(c, string(fullURL))
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	exShortURL := config.ExpShortURL(shortURL)
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.String(http.StatusCreated, exShortURL)
}

func (h ServerHandler) GetFullURL(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	shortURL := c.Param("id")
	fullURL, err := h.Storage.GetFullURL(c, shortURL)
	if err != nil {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	if !strings.HasPrefix(fullURL, config.HTTP) {
		fullURL = config.HTTP + strings.TrimPrefix(fullURL, "//")
	}
	c.Redirect(http.StatusTemporaryRedirect, fullURL)

}

func (h ServerHandler) GetShortURL(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	var full repository.FullURL
	err = json.Unmarshal(reqBody, &full)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	var sURL repository.ShortURL
	var responseStatus int
	sURL.Short, err = h.Storage.GetShortURL(c, full.Full)
	if err != nil {
		fromInsert, err := h.Storage.InsertURL(c, full.Full)
		if err != nil {
			c.String(http.StatusBadRequest, "")
			return
		}
		sURL.Short = fromInsert
		responseStatus = http.StatusCreated
	} else {
		responseStatus = http.StatusOK
	}
	sURL.Short = config.ExpShortURL(sURL.Short)
	respBody, err := json.Marshal(sURL)
	if err != nil {
		c.String(http.StatusNotFound, "")
	}
	c.String(responseStatus, string(respBody))
}

func (h ServerHandler) ResponseBadRequest(c *gin.Context) {
	c.String(http.StatusBadRequest, "")
}
