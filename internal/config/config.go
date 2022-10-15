package config

import (
	"github.com/caarlos0/env"
	"os"
	"strings"
)

const (
	HostPort string = "8080"
	HostAddr string = "localhost"
	HTTP     string = "http://"
)

var Cnf Config

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func NewConfig(options ...Option) *Config {
	conf := Config{
		ServerAddress: HostAddr + ":" + HostPort,
		BaseURL:       HostAddr + ":" + HostPort,
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
	}
}

func WithServerAddress(hostAddr, hostPort string) Option {
	return func(s *Config) {
		s.ServerAddress = hostAddr + ":" + hostPort
	}
}

func WithBaseURL(hostAddr, hostPort string) Option {
	return func(s *Config) {
		s.BaseURL = hostAddr + ":" + hostPort
	}
}

func (c *Config) BasePort() string {
	part := strings.Split(c.BaseURL, ":")
	cnt := len(part)
	if cnt > 1 {
		return part[cnt-1]
	} else {
		return HostPort
	}
}

func (c *Config) HostAddr() string {
	part := strings.Split(c.BaseURL, ":")
	if strings.HasPrefix(part[0], HTTP) {
		str := strings.TrimSuffix(part[0], HTTP)
		return str
	}
	return part[0]
}

func ExpShortURL(shortURL string) string {
	x, err := os.LookupEnv("BASE_URL")
	if err {
		return x + "/" + shortURL
	}
	return HTTP + HostAddr + ":" + Cnf.BasePort() + "/" + shortURL
}
