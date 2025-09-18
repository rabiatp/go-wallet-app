package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort     string
	DatabaseURL string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/wallet_db?sslmode=disable")

	viper.AutomaticEnv()
	cfg := &Config{
		AppPort:     viper.GetString("APP_PORT"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
	}
	log.Printf("config loaded at %s", time.Now().Format(time.RFC3339))
	return cfg
}
