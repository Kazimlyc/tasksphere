package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

func Load() *Config {
	loadEnv()

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Port:        os.Getenv("PORT"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is missing")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWTSecret is missing")
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	return cfg
}

func loadEnv() {
	// Try current directory first, fall back to the repo root when backend is executed from its own folder.
	loaders := []string{".env", filepath.Join("..", ".env")}
	for _, path := range loaders {
		if err := godotenv.Load(path); err == nil {
			return
		}
	}
}
