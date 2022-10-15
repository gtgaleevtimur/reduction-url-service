package app

import (
	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	hd "github.com/gtgaleevtimur/reduction-url-service/internal/handlers"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
	"log"
)

func Run() {
	conf := config.NewConfig(config.WithParseEnv())
	storage := repository.NewStorage()
	log.Fatal(hd.NewRouter(storage).Run(conf.ServerAddress))
}
