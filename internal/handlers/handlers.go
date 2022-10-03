package handlers

import (
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
	r := gin.Default()
	r.POST("/", controller.CreateShortURL)
	r.GET("/:id", controller.GetFullURL)
	r.NoRoute(controller.ResponseBadRequest)
	return r
}

type ServerHandler struct {
	Storage *repository.Storage
}

func NewServerHandler(s *repository.Storage) *ServerHandler {
	return &ServerHandler{Storage: s}
}

func (h ServerHandler) CreateShortURL(c *gin.Context) {
	fullURL, err := ioutil.ReadAll(c.Request.Body)
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

func (h ServerHandler) ResponseBadRequest(c *gin.Context) {
	c.String(http.StatusBadRequest, "")
}
