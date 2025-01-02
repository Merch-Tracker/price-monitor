package config

import "os"

type Config struct {
	AppConfig  AppConfig
	RepoConfig RepoConfig
}

type AppConfig struct {
	Host     string
	Port     string
	LogLevel string
}

type RepoConfig struct {
	Host string
	Port string
}

func NewConfig() *Config {
	return &Config{
		AppConfig: AppConfig{
			Host:     getEnv("APP_HOST", "0.0.0.0"),
			Port:     getEnv("APP_PORT", "9050"),
			LogLevel: getEnv("APP_LOG_LEVEL", "info"),
		},
		RepoConfig: RepoConfig{
			Host: getEnv("REPO_HOST", "0.0.0.0"),
			Port: getEnv("REPO_PORT", "9010"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
