package configs

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Environment  string        `mapstructure:"ENVIRONMENT"`
	DBHost       string        `mapstructure:"DB_HOST"`
	DBPort       string        `mapstructure:"DB_PORT"`
	DBUser       string        `mapstructure:"DB_USER"`
	DBPassword   string        `mapstructure:"DB_PASSWORD"`
	DBName       string        `mapstructure:"DB_NAME"`
	ServerPort   string        `mapstructure:"SERVER_PORT"`
	JWTSecret    string        `mapstructure:"JWT_SECRET"`
	TokenExpiry  time.Duration `mapstructure:"TOKEN_EXPIRY"`
	AllowOrigins string        `mapstructure:"ALLOW_ORIGINS"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Info().Msg("No .env file found or error loading it. Using environment variables.")
	}
	config := &Config{
		Environment:  getEnv("ENVIRONMENT", "development"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "postgres"),
		DBName:       getEnv("DB_NAME", "wallet"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		TokenExpiry:  time.Duration(getEnvAsInt("TOKEN_EXPIRY", 24)) * time.Hour,
		AllowOrigins: getEnv("ALLOW_ORIGINS", "*"),
	}
	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
