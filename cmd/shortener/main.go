package main

import (
	"fmt"
	"github.com/gtgaleevtimur/reduction-url-service/internal/app"
)

// Константы версий приложения.
var (
	buildVersion = "1.0.0"
	buildDate    = "24.01.2023"
	buildCommit  = "1.0.0"
)

func main() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
	// Через единственный вход запускаем приложение.
	app.Run()
}
