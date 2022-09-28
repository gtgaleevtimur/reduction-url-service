package main

import "github.com/gtgaleevtimur/reduction-url-service/internal/app"

const addrServ = "http://localhost:8080"

func main() {
	app.Run(addrServ)
}
