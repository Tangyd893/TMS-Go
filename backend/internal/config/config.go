package config

import (
	"os"
	"strconv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	MinIO    MinIOConfig
	RabbitMQ RabbitMQConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	Name         string
	User         string
	Password     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessSecret    string
	RefreshSecret   string
	AccessTokenTTL  string
	RefreshTokenTTL string
	Issuer          string
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type RabbitMQConfig struct {
	URL string
}

func Load() Config {
	return Config{
		App: AppConfig{
			Name: envOrDefault("APP_NAME", "tms-api"),
			Env:  envOrDefault("APP_ENV", "dev"),
			Port: envOrDefault("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:         envOrDefault("DB_HOST", "localhost"),
			Port:         envOrDefault("DB_PORT", "5432"),
			Name:         envOrDefault("DB_NAME", "tms"),
			User:         envOrDefault("DB_USER", "tms"),
			Password:     envOrDefault("DB_PASSWORD", "change_me"),
			SSLMode:      envOrDefault("DB_SSLMODE", "disable"),
			MaxOpenConns: envOrDefaultInt("DB_MAX_OPEN_CONNS", 50),
			MaxIdleConns: envOrDefaultInt("DB_MAX_IDLE_CONNS", 10),
		},
		Redis: RedisConfig{
			Host:     envOrDefault("REDIS_HOST", "localhost"),
			Port:     envOrDefault("REDIS_PORT", "6379"),
			Password: envOrDefault("REDIS_PASSWORD", ""),
			DB:       envOrDefaultInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			AccessSecret:    envOrDefault("JWT_ACCESS_SECRET", "tms-access-secret-change-me"),
			RefreshSecret:   envOrDefault("JWT_REFRESH_SECRET", "tms-refresh-secret-change-me"),
			AccessTokenTTL:  envOrDefault("JWT_ACCESS_TTL", "2h"),
			RefreshTokenTTL: envOrDefault("JWT_REFRESH_TTL", "168h"),
			Issuer:          envOrDefault("JWT_ISSUER", "tms"),
		},
		MinIO: MinIOConfig{
			Endpoint:  envOrDefault("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: envOrDefault("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: envOrDefault("MINIO_SECRET_KEY", "change_me"),
			Bucket:    envOrDefault("MINIO_BUCKET", "tms-dev"),
			UseSSL:    envOrDefault("MINIO_USE_SSL", "false") == "true",
		},
		RabbitMQ: RabbitMQConfig{
			URL: envOrDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
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

func envOrDefaultInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return n
}
