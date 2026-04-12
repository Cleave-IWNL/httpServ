package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL    string
	MigrationsPath string
	ServerPort     string
	LogLevel       string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	dbUrl := viper.GetString("DATABASE_URL")

	if dbUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return &Config{
		DatabaseURL:    dbUrl,
		MigrationsPath: viper.GetString("MIGRATIONS_PATH"),
		ServerPort:     viper.GetString("SERVER_PORT"),
		LogLevel:       viper.GetString("LOG_LEVEL"),
	}, nil
}
