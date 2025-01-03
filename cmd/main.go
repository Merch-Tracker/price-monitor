package main

import (
	"parsing-service/config"
	"parsing-service/internal/app"
	"parsing-service/internal/logging"
)

func main() {
	c := config.NewConfig()

	logging.LogSetup(c.LogLevel)

	appl := app.New(c)

	appl.Run()
}
