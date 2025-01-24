package config

import "os"

type Config struct {
	DB PostgresConfig
}

type PostgresConfig struct {
	User     string
	Password string
	URL      string
	Port     string
	DBName   string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DB: PostgresConfig{
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			URL:      os.Getenv("POSTGRES_URL"),
			Port:     os.Getenv("POSTGRES_PORT"),
			DBName:   os.Getenv("POSTGRES_DBNAME"),
		}}
	return cfg, nil
}
