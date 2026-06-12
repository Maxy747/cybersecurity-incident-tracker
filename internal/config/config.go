package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ServerPort string
}

// Load reads the .env file (if present) and returns a populated Config.
func Load() *Config {
	// Best-effort: ignore error if .env doesn't exist (e.g. in Docker with env vars set directly)
	if err := godotenv.Load(); err != nil {
		log.Println("[config] .env file not found, reading from environment")
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "cyberguard"),
		DBPassword: getEnv("DB_PASSWORD", "cyberguard_secret"),
		DBName:     getEnv("DB_NAME", "cyberguard_db"),
		JWTSecret:  getEnv("JWT_SECRET", "change_me_in_production"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
	return cfg
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
