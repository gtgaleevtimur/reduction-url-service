package config

import (
	"strings"
)

const (
	HostPort string = "8080"
	HostAddr string = "localhost"
	HTTP     string = "http://"
)

var Cnf Config

type Config struct {
	ServerAddress string
	BaseURL       string
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

type Option func(s *Config)

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

func ExpShortURL(shortURL string) string {
	return HTTP + HostAddr + ":" + Cnf.BasePort() + "/" + shortURL
}
