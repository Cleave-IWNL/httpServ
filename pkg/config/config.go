package config

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL        string
	MigrationsPath     string
	ServerPort         string
	LogLevel           string
	ExchangeAPIURL     string
	ExchangeAPIKey     string
	KafkaBrokers       []string
	KafkaPaymentsTopic string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	kafkaBrokers := strings.Split(viper.GetString("KAFKA_BROKERS"), ",")

	for i, b := range kafkaBrokers {
		kafkaBrokers[i] = strings.TrimSpace(b)
	}

	if slices.Contains(kafkaBrokers, "") {
		return nil, fmt.Errorf("KAFKA_BROKERS is required")
	}

	kafkaPaymentsTopic := viper.GetString("KAFKA_PAYMENTS_TOPIC")

	if kafkaPaymentsTopic == "" {
		return nil, fmt.Errorf("KAFKA_PAYMENTS_TOPIC is required")
	}

	dbUrl := viper.GetString("DATABASE_URL")

	if dbUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	exchangeKey := viper.GetString("EXCHANGE_API_KEY")

	if exchangeKey == "" {
		return nil, fmt.Errorf("EXCHANGE_API_KEY is required")
	}

	exchangeURL := viper.GetString("EXCHANGE_API_URL")

	if exchangeURL == "" {
		exchangeURL = "https://v6.exchangerate-api.com"
	}

	return &Config{
		DatabaseURL:        dbUrl,
		MigrationsPath:     viper.GetString("MIGRATIONS_PATH"),
		ServerPort:         viper.GetString("SERVER_PORT"),
		LogLevel:           viper.GetString("LOG_LEVEL"),
		ExchangeAPIURL:     exchangeURL,
		ExchangeAPIKey:     exchangeKey,
		KafkaBrokers:       kafkaBrokers,
		KafkaPaymentsTopic: kafkaPaymentsTopic,
	}, nil
}
