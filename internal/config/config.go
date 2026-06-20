package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	AppVersion    string
	DatabasePath  string
	BCryptCost    int
	JWTTimeToLive int
	JWTSecretKey  string
	Domain        string
	ShortIDLength int
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

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = ":8080"

		slog.Warn("Environment Variable not found, using default",
			slog.String("DOMAIN", domain))
	}

	appVersion := os.Getenv("APP_VERSION")
	if appVersion == "" {
		appVersion = "undefined"

		slog.Warn("Environment Variable not found, using default",
			slog.String("APP_VERSION", appVersion))
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/forma.db"

		slog.Warn("Environment Variable not found, using default",
			slog.String("DB_PATH", dbPath))
	}

	var cost int
	costStr := os.Getenv("BCRYPT_COST")
	if costStr == "" {
		costStr = "12"
		cost, _ = strconv.Atoi(costStr)

		slog.Warn("Environment Variable not found, using default",
			slog.String("BCRYPT_COST", costStr))
	} else {
		cost, _ = strconv.Atoi(costStr)
	}

	var ttl int
	ttlStr := os.Getenv("JWT_TTL")
	if ttlStr == "" {
		ttlStr = "4380"
		ttl, _ = strconv.Atoi(ttlStr)

		slog.Warn("Environment Variable not found, using default",
			slog.String("JWT_TTL", ttlStr))
	} else {
		ttl, _ = strconv.Atoi(ttlStr)
	}

	var length int
	lengthStr := os.Getenv("SHORT_ID_LENGTH")
	if lengthStr == "" {
		lengthStr = "8"
		length, _ = strconv.Atoi(lengthStr)

		slog.Warn("Environment Variable not found, using default",
			slog.String("SHORT_ID_LENGTH", lengthStr))
	} else {
		length, _ = strconv.Atoi(lengthStr)
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		slog.Error("failed to load environment variable jwt secret key")
	}

	return &Config{
		ServerPort:    port,
		AppVersion:    appVersion,
		DatabasePath:  dbPath,
		BCryptCost:    cost,
		JWTTimeToLive: ttl,
		JWTSecretKey:  secretKey,
		Domain:        domain,
		ShortIDLength: length,
	}
}
