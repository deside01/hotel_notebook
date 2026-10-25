package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   string
	DB_URL string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("empty env")
	}

	cfg := &Config{
		Port:   os.Getenv("PORT"),
		DB_URL: os.Getenv("DB_URL"),
	}

	return cfg, nil
}
