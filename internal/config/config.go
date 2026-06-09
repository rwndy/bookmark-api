package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBDsn           string
	JWTSecret       string
	Port            string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func Load() *Config {
	return &Config{
		DBDsn:           getEnv("DB_DSN", ""),
		JWTSecret:       getEnv("JWT_SECRET", "default-secret"),
		Port:            getEnv("PORT", "3000"),
		AccessTokenTTL:  getEnvDuration("ACCESS_TOKEN_TTL_MIN", 15) * time.Minute,
		RefreshTokenTTL: getEnvDuration("REFRESH_TOKEN_TTL_HOUR", 24*7) * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvDuration(key string, fallback int) time.Duration {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return time.Duration(n)
		}
	}
	return time.Duration(fallback)
}
