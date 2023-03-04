package config

import (
	"encoding/json"
	"flag"
	"os"
	"strings"
	"sync"

	"github.com/caarlos0/env"
)

// Настройки сервера по умолчанию.
const (
	HostPort string = "8080"      // порт хоста по дефолту.
	HostAddr string = "localhost" // адрес хоста по дефолту.
	HTTP     string = "http://"   // префикс адреса по дефолту.
)

var (
	config *Config
	once   sync.Once
)

// Config - структура конфигурационного файла приложения.
type Config struct {
	ServerAddress string `json:"server_address" env:"SERVER_ADDRESS"`
	BaseURL       string `json:"base_url" env:"BASE_URL"`
	StoragePath   string `json:"file_storage_path" env:"FILE_STORAGE_PATH"`
	DatabaseDSN   string `json:"database_dsn" env:"DATABASE_DSN"`
	EnableHTTPS   bool   `json:"enable_https" env:"ENABLE_HTTPS"`
	Config        string `env:"CONFIG"`
	TrustedSubnet string `json:"trusted_subnet" env:"TRUSTED_SUBNET"`
	EnableGRPC    bool   `json:"enable_grpc" env:"ENABLE_GRPC"`
}

// NewConfig - конструктор конфигурационного файла.
func NewConfig(options ...Option) *Config {
	once.Do(
		func() {
			config = &Config{
				ServerAddress: HostAddr + ":" + HostPort,
				BaseURL:       HostAddr + ":" + HostPort,
				StoragePath:   "",
				DatabaseDSN:   "",
				Config:        "",
				TrustedSubnet: "",
				EnableGRPC:    false,
			}

			// если в аргументах получили Options, то применяем их к Config.
			for _, opt := range options {
				opt(config)
			}
			configDataJSON, err := os.ReadFile(config.Config)
			if err != nil {
				return
			}
			var configJSON Config
			if err = json.Unmarshal(configDataJSON, &configJSON); err != nil {
				return
			}
			if config.ServerAddress == "" {
				config.ServerAddress = configJSON.ServerAddress
			}
			if config.BaseURL == "" {
				config.BaseURL = configJSON.ServerAddress
			}
			if config.StoragePath == "" {
				config.StoragePath = configJSON.StoragePath
			}
			if config.DatabaseDSN == "" {
				config.DatabaseDSN = configJSON.StoragePath
			}
			if config.TrustedSubnet == "" {
				config.TrustedSubnet = configJSON.TrustedSubnet
			}
			if !config.EnableGRPC {
				config.EnableGRPC = configJSON.EnableGRPC
			}
			if !config.EnableHTTPS {
				config.EnableHTTPS = configJSON.EnableHTTPS
			}
		})

	return config
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
	flag.StringVar(&c.TrustedSubnet, "t", c.TrustedSubnet, "TRUSTED_SUBNET")
	flag.BoolVar(&c.EnableGRPC, "g", c.EnableGRPC, "ENABLE_GRPC")
	flag.Parse()
}

// ExpShortURL - хэлпер, формирующий сокращенный URL (http+hostAddr+hostport+hash).
func (c *Config) ExpShortURL(shortURL string) string {
	if strings.HasPrefix(c.BaseURL, HTTP) {
		return c.BaseURL + "/" + shortURL
	}
	return HTTP + HostAddr + ":" + HostPort + "/" + shortURL
}
