package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                 string
	Port                string
	DatabaseURL         string
	JWTAccessSecret     string
	JWTRefreshSecret    string
	AccessTTL           time.Duration
	RefreshTTL          time.Duration
	AdminBootstrapToken string
	AdminEmail          string
	AdminPassword       string
	AllowedOrigins      []string
}

func Load() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	cfg := &Config{
		Env:                 getEnv("ENV", "development"),
		Port:                getEnv("PORT", "8080"),
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		JWTAccessSecret:     getEnv("JWT_ACCESS_SECRET", ""),
		JWTRefreshSecret:    getEnv("JWT_REFRESH_SECRET", ""),
		AdminBootstrapToken: getEnv("ADMIN_BOOTSTRAP_TOKEN", ""),
		AdminEmail:          getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword:       getEnv("ADMIN_PASSWORD", ""),
		AllowedOrigins:      getEnvSlice("ALLOWED_ORIGINS", []string{"*"}),
	}

	// Parse durations
	var err error
	cfg.AccessTTL, err = time.ParseDuration(getEnv("ACCESS_TTL", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid ACCESS_TTL: %w", err)
	}

	cfg.RefreshTTL, err = time.ParseDuration(getEnv("REFRESH_TTL", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TTL: %w", err)
	}

	// Validate required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWTAccessSecret == "" {
		return nil, fmt.Errorf("JWT_ACCESS_SECRET is required")
	}
	if cfg.JWTRefreshSecret == "" {
		return nil, fmt.Errorf("JWT_REFRESH_SECRET is required")
	}

	log.Printf("Config loaded: ENV=%s, PORT=%s", cfg.Env, cfg.Port)
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// Simple split by comma
	var result []string
	current := ""
	for _, char := range value {
		if char == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}

	return result
}
