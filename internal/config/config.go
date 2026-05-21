package config

import "os"

type Config struct {
	DBDsn     string
	JWTSecret string
	Port      string
}

func Load() *Config {
	return &Config{
		DBDsn:     getEnv("DB_DSN", ""),
		JWTSecret: getEnv("JWT_SECRET", "default-secret"),
		Port:      getEnv("PORT", "3000"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}