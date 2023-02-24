package config

import (
	"encoding/json"
	"flag"
	"os"
	"strings"

	"github.com/caarlos0/env"
)

// Настройки сервера по умолчанию.
const (
	HostPort string = "8080"      // порт хоста по дефолту.
	HostAddr string = "localhost" // адрес хоста по дефолту.
	HTTP     string = "http://"   // префикс адреса по дефолту.
)

// Config - структура конфигурационного файла приложения.
type Config struct {
	ServerAddress string `json:"server_address" env:"SERVER_ADDRESS"`
	BaseURL       string `json:"base_url" env:"BASE_URL"`
	StoragePath   string `json:"file_storage_path" env:"FILE_STORAGE_PATH"`
	DatabaseDSN   string `json:"database_dsn" env:"DATABASE_DSN"`
	EnableHTTPS   bool   `json:"enable_https" env:"ENABLE_HTTPS"`
	Config        string `env:"CONFIG"`
}

// NewConfig - конструктор конфигурационного файла.
func NewConfig(options ...Option) *Config {
	conf := Config{
		ServerAddress: HostAddr + ":" + HostPort,
		BaseURL:       HostAddr + ":" + HostPort,
		StoragePath:   "",
		DatabaseDSN:   "",
		Config:        "",
	}

	// если в аргументах получили Options, то применяем их к Config.
	for _, opt := range options {
		opt(&conf)
	}
	configDataJSON, err := os.ReadFile(conf.Config)
	if err != nil {
		return &conf
	}
	var configJSON Config
	if err = json.Unmarshal(configDataJSON, &configJSON); err != nil {
		return &conf
	}
	if conf.ServerAddress == "" {
		conf.ServerAddress = configJSON.ServerAddress
	}
	if conf.BaseURL == "" {
		conf.BaseURL = configJSON.ServerAddress
	}
	if conf.StoragePath == "" {
		conf.StoragePath = configJSON.StoragePath
	}
	if conf.DatabaseDSN == "" {
		conf.DatabaseDSN = configJSON.StoragePath
	}
	if !conf.EnableHTTPS {
		conf.EnableHTTPS = configJSON.EnableHTTPS
	}
	return &conf
}

// Option - функция применяемая к Config для его заполнения.
type Option func(*Config)

// WithParseEnv - парсит из окружения/флагов, изменяет Config.
func WithParseEnv() Option {
	return func(c *Config) {
		env.Parse(c)
		c.ParseFlags()
	}
}

// ParseFlags - парсит флаги.
func (c *Config) ParseFlags() {
	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "SERVER_ADDRESS")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "BASE_URL")
	flag.StringVar(&c.StoragePath, "f", c.StoragePath, "FILE_STORAGE_PATH")
	flag.StringVar(&c.DatabaseDSN, "d", c.DatabaseDSN, "DATABASE_DSN")
	flag.BoolVar(&c.EnableHTTPS, "s", c.EnableHTTPS, "ENABLE_HTTPS")
	flag.StringVar(&c.Config, "c", c.Config, "config JSON file")
	flag.Parse()
}

// ExpShortURL - хэлпер, формирующий сокращенный URL (http+hostAddr+hostport+hash).
func (c *Config) ExpShortURL(shortURL string) string {
	if strings.HasPrefix(c.BaseURL, HTTP) {
		return c.BaseURL + "/" + shortURL
	}
	return HTTP + HostAddr + ":" + HostPort + "/" + shortURL
}
