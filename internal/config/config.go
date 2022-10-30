package config

import (
	"flag"
	"strings"

	"github.com/caarlos0/env"
)

const (
	HostPort string = "8080"
	HostAddr string = "localhost"
	HTTP     string = "http://"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	StoragePath   string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
}

func NewConfig(options ...Option) *Config {
	conf := Config{
		ServerAddress: HostAddr + ":" + HostPort,
		BaseURL:       HostAddr + ":" + HostPort,
		StoragePath:   "",
		DatabaseDSN:   "",
	}

	for _, opt := range options {
		opt(&conf)
	}
	return &conf
}

type Option func(*Config)

func WithParseEnv() Option {
	return func(c *Config) {
		env.Parse(c)
		c.ParseFlags()
	}
}

func (c *Config) ParseFlags() {
	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "SERVER_ADDRESS")
	flag.StringVar(&c.BaseURL, "b", c.BaseURL, "BASE_URL")
	flag.StringVar(&c.StoragePath, "f", c.StoragePath, "FILE_STORAGE_PATH")
	flag.StringVar(&c.DatabaseDSN, "d", c.DatabaseDSN, "DATABASE_DSN")
	flag.Parse()
}

func (c *Config) ExpShortURL(shortURL string) string {
	if strings.HasPrefix(c.BaseURL, HTTP) {
		return c.BaseURL + "/" + shortURL
	}
	return HTTP + HostAddr + ":" + HostPort + "/" + shortURL
}
