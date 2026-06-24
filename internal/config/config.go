package config

import "os"

type AppConfig struct {
	Host      string
	Port      string
	JWTSecret string
	APIKey    string
	StripeKey string
	Mode      string
	Node      string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type CacheConfig struct {
	Host     string
	Password string
}

type Config struct {
	App   AppConfig
	DB    DBConfig
	Cache CacheConfig
}

func LoadConfig() Config {
	return Config{
		App: AppConfig{
			JWTSecret: getEnv("JWT_SECRET", "default_jwt_secret"),
			APIKey:    getEnv("API_KEY", "default_api_key"),
			StripeKey: getEnv("STRIPE_KEY", "default_stripe_key"),
			Host:      getEnv("API_HOST", "localhost"),
			Port:      getEnv("API_PORT", "4000"),
			Mode:      getEnv("APP_MODE", "release"),
			Node:      getEnv("NODE_ID", "1"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "app_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Cache: CacheConfig{
			Host:     getEnv("CACHE_URI", "localhost:6379"),
			Password: getEnv("CACHE_PASSWORD", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
