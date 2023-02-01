package main

import (
	"fmt"

	"github.com/gtgaleevtimur/reduction-url-service/internal/app"
)

// Константы версий приложения.
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
	// Через единственный вход запускаем приложение.
	app.Run()
}
