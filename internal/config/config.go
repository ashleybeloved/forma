package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort         string
	AppVersion         string
	HTTPS              bool
	DatabasePath       string
	GeoIPDatabasePath  string
	BCryptCost         int
	JWTTimeToLive      int
	JWTSecretKey       string
	Domain             string
	ShortIDLength      int
	PasswordMinSymbols int
	PasswordMaxSymbols int
	UsernameMinSymbols int
	UsernameMaxSymbols int
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("failed to initialize .env file")
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		slog.Error("failed to load environment variable jwt secret key, please input secret key to env!")
		os.Exit(1)
	}

	return &Config{
		ServerPort:         getEnvString("SERVER_PORT", ":8080"),
		Domain:             getEnvString("DOMAIN", ":8080"),
		AppVersion:         getEnvString("APP_VERSION", "undefined"),
		DatabasePath:       getEnvString("DB_PATH", "./data/forma.db"),
		GeoIPDatabasePath:  getEnvString("GEOIP_DB_PATH", "./data/GeoLite2-Country.mmdb"),
		JWTSecretKey:       secretKey,
		HTTPS:              getEnvBool("HTTPS", true),
		BCryptCost:         getEnvInt("BCRYPT_COST", 12),
		JWTTimeToLive:      getEnvInt("JWT_TTL", 4380),
		ShortIDLength:      getEnvInt("SHORT_ID_LENGTH", 8),
		PasswordMinSymbols: getEnvInt("PASSWORD_MIN", 4),
		PasswordMaxSymbols: getEnvInt("PASSWORD_MAX", 128),
		UsernameMinSymbols: getEnvInt("USERNAME_MIN", 4),
		UsernameMaxSymbols: getEnvInt("USERNAME_MAX", 32),
	}
}

func getEnvString(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		slog.Warn("Environment Variable not found, using default", slog.String(key, defaultValue))
		return defaultValue
	}
	return val
}

func getEnvBool(key string, defaultValue bool) bool {
	valStr := os.Getenv(key)
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		slog.Warn("Environment Variable not found or invalid, using default", slog.String(key, strconv.FormatBool(defaultValue)))
		return defaultValue
	}
	return val
}

func getEnvInt(key string, defaultValue int) int {
	valStr := os.Getenv(key)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		slog.Warn("Environment Variable not found or invalid, using default", slog.String(key, strconv.Itoa(defaultValue)))
		return defaultValue
	}
	return val
}
