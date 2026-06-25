package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseURL string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("JWT_SECRET not set, using fallback key")
		jwtSecret = "default_secret_key"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL not set, using fallback connection string")
		dbURL = "postgres://postgres:password@localhost:5432/tickets?sslmode=disable"
	}

	return &Config{
		Port:        port,
		JWTSecret:   jwtSecret,
		DatabaseURL: dbURL,
	}
}
