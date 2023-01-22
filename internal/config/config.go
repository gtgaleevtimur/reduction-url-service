package config

import (
	"flag"
	"strings"

	"github.com/caarlos0/env"
)

const (
	// HostPort - порт хоста по дефолту.
	HostPort string = "8080"
	// HostAddr - адрес хоста по дефолту.
	HostAddr string = "localhost"
	// HTTP - префикс адреса по дефолту.
	HTTP string = "http://"
)

// Config - структура конфигурационного файла приложения.
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	StoragePath   string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
}

// NewConfig - конструктор конфигурационного файла.
func NewConfig(options ...Option) *Config {
	conf := Config{
		ServerAddress: HostAddr + ":" + HostPort,
		BaseURL:       HostAddr + ":" + HostPort,
		StoragePath:   "",
		DatabaseDSN:   "",
	}

	// если в аргументах получили Options, то применяем их к Config.
	for _, opt := range options {
		opt(&conf)
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
	flag.Parse()
}

// ExpShortURL - хэлпер, формирующий сокращенный URL (http+hostAddr+hostport+hash).
func (c *Config) ExpShortURL(shortURL string) string {
	if strings.HasPrefix(c.BaseURL, HTTP) {
		return c.BaseURL + "/" + shortURL
	}
	return HTTP + HostAddr + ":" + HostPort + "/" + shortURL
}
