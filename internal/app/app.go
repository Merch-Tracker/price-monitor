package app

import (
	log "github.com/sirupsen/logrus"
	"parsing-service/config"
	"runtime"
)

type App struct {
	Address string
	NumCPUs int
}

func New(c *config.Config) *App {
	numCPUs := c.AppConfig.NumCPUs
	if numCPUs < 1 {
		numCPUs = runtime.NumCPU()
	}

	return &App{
		Address: c.AppConfig.Host + ":" + c.AppConfig.Port,
		NumCPUs: numCPUs,
	}
}

func (app *App) Start() {
	log.Info("Application start")
	log.WithFields(log.Fields{
		"Address":        app.Address,
		"Number of CPUs": app.NumCPUs,
	}).Debug("App settings")
}
