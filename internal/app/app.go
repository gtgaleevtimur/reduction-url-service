package app

import (
	"log"

	"github.com/gtgaleevtimur/reduction-url-service/internal/config"
	hd "github.com/gtgaleevtimur/reduction-url-service/internal/handlers"
	"github.com/gtgaleevtimur/reduction-url-service/internal/repository"
)

func Run() {
	conf := config.NewConfig(config.WithParseEnv())
	storage := repository.NewStorage(conf)
	log.Fatal(hd.NewRouter(storage, conf).Run(conf.ServerAddress))
}
