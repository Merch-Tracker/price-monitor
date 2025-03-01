package config

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type Config struct {
	Host        string
	ClientPort  string
	ServerPort  string
	LogLevel    string
	NumCPUs     int
	CheckPeriod int
}

func NewConfig() *Config {
	return &Config{
		Host:        getEnv("APP_HOST", "0.0.0.0"),
		ClientPort:  getEnv("APP_PORT", "9050"),
		ServerPort:  getEnv("APP_SERVER_PORT", "9060"),
		LogLevel:    getEnv("APP_LOG_LEVEL", "debug"),
		NumCPUs:     getEnvInt("APP_NUMCPUS", -1),
		CheckPeriod: getEnvInt("APP_CHECK_PERIOD", 6),
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
