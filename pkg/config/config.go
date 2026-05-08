package config

import (
	"fmt"
	"slices"
	"strings"
	"time"

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
	OutboxWorkerCount  int
	OutboxBatchSize    int
	OutboxPollInterval time.Duration
	OutboxMaxAttempts  int
	OutboxBaseBackoff  time.Duration
	OutboxMaxBackoff   time.Duration
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

	outboxWorkers := viper.GetInt("OUTBOX_WORKERS")
	if outboxWorkers <= 0 {
		return nil, fmt.Errorf("OUTBOX_WORKERS must be > 0")
	}

	outboxBatchSize := viper.GetInt("OUTBOX_BATCH_SIZE")
	if outboxBatchSize <= 0 {
		return nil, fmt.Errorf("OUTBOX_BATCH_SIZE must be > 0")
	}

	outboxPollInterval := viper.GetDuration("OUTBOX_POLL_INTERVAL")
	if outboxPollInterval <= 0 {
		return nil, fmt.Errorf("OUTBOX_POLL_INTERVAL must be > 0")
	}

	outboxMaxAttempts := viper.GetInt("OUTBOX_MAX_ATTEMPTS")
	if outboxMaxAttempts <= 0 {
		return nil, fmt.Errorf("OUTBOX_MAX_ATTEMPTS must be > 0")
	}

	outboxBaseBackoff := viper.GetDuration("OUTBOX_BASE_BACKOFF")
	if outboxBaseBackoff <= 0 {
		return nil, fmt.Errorf("OUTBOX_BASE_BACKOFF must be > 0")
	}

	outboxMaxBackoff := viper.GetDuration("OUTBOX_MAX_BACKOFF")
	if outboxMaxBackoff <= 0 {
		return nil, fmt.Errorf("OUTBOX_MAX_BACKOFF must be > 0")
	}

	if outboxBaseBackoff > outboxMaxBackoff {
		return nil, fmt.Errorf("OUTBOX_BASE_BACKOFF must be <= OUTBOX_MAX_BACKOFF")
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
		OutboxWorkerCount:  outboxWorkers,
		OutboxBatchSize:    outboxBatchSize,
		OutboxPollInterval: outboxPollInterval,
		OutboxMaxAttempts:  outboxMaxAttempts,
		OutboxBaseBackoff:  outboxBaseBackoff,
		OutboxMaxBackoff:   outboxMaxBackoff,
	}, nil
}
