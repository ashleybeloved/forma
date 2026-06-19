package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	AppVersion string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("failed to initialize .env file")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":8080"

		slog.Warn("Environment Variable not found, using default",
			slog.String("SERVER_PORT", port))
	}

	appVersion := os.Getenv("APP_VERSION")
	if appVersion == "" {
		appVersion = "undefined"

		slog.Warn("Environment Variable not found, using default",
			slog.String("APP_VERSION", appVersion))
	}

	return &Config{
		ServerPort: port,
		AppVersion: appVersion,
	}
}
