package main

import (
	"parsing-service/config"
	"parsing-service/internal/app"
	"parsing-service/internal/logging"
)

func main() {
	c := config.NewConfig()

	logging.LogSetup(c.AppConfig.LogLevel)

	appl := app.New(c)

	appl.Start()
}
