package config

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type Config struct {
	Host     string
	Port     string
	LogLevel string
	NumCPUs  int
}

func NewConfig() *Config {
	return &Config{
		Host:     getEnv("APP_HOST", "0.0.0.0"),
		Port:     getEnv("APP_PORT", "9050"),
		LogLevel: getEnv("APP_LOG_LEVEL", "debug"),
		NumCPUs:  getEnvInt("APP_NUMCPUs", -1),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		num, err := strconv.Atoi(value)
		if err != nil {
			log.WithField("default", -1).Warn("Config | Can't parse value as int")
			return fallback
		}
		return num
	}
	return fallback
}
