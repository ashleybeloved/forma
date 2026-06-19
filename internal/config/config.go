package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("failed to initialize .env file")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":8080"

		slog.Warn("SERVER_PORT not found in .env, using default",
			slog.String("SERVER_PORT", ":8080"))
	}

	return &Config{
		ServerPort: port,
	}
}
