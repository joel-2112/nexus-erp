package config

import (
	"os"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
}

func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8081"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "5433"),
		DBUser:     getEnv("DB_USER", "nexus_user"),
		DBPassword: getEnv("DB_PASSWORD", "nexus_password"),
		DBName: getEnv("DB_NAME", "nexus_auth_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", "nexus-dev-secret"),
	}
}

func (c *Config) DatabaseURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword +
		"@" + c.DBHost + ":" + c.DBPort +
		"/" + c.DBName + "?sslmode=" + c.DBSSLMode
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}