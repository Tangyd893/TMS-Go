package config

import (
	"os"
)

type Config struct {
	App AppConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

func Load() Config {
	return Config{
		App: AppConfig{
			Name: envOrDefault("APP_NAME", "tms-api"),
			Env:  envOrDefault("APP_ENV", "dev"),
			Port: envOrDefault("APP_PORT", "8080"),
		},
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
