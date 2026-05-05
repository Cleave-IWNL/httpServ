package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"httpServ/internal/client/exchangerate"
	"httpServ/internal/handler"
	"httpServ/internal/repository"
	"httpServ/internal/service"
	"httpServ/pkg/config"
	"httpServ/pkg/db"
	"httpServ/pkg/httpclient"
	kafkaclient "httpServ/pkg/kafka"
	"httpServ/pkg/logger"
	"httpServ/worker"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	zapLog, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer zapLog.Sync()

	database, err := db.New(db.DatabaseConfig{
		URL:            cfg.DatabaseURL,
		MigrationsPath: cfg.MigrationsPath,
	})
	if err != nil {
		zapLog.Fatal("failed to init db", zap.Error(err))
	}
	defer database.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	baseHTTP := httpclient.NewDefaultClient(5 * time.Second)
	retryHTTP := httpclient.NewRetryClient(baseHTTP, httpclient.RetryConfig{
		MaxRetries: 3,
	}, zapLog)
	loggedHTTP := httpclient.NewLoggingClient(retryHTTP, zapLog)

	rateProvider := exchangerate.New(loggedHTTP, cfg.ExchangeAPIURL, cfg.ExchangeAPIKey)

	producer, err := kafkaclient.New(cfg.KafkaBrokers, cfg.KafkaPaymentsTopic, zapLog)
	if err != nil {
		zapLog.Fatal("failed to init kafka producer", zap.Error(err))
	}
	defer producer.Close()

	outboxRepo := repository.NewOutboxPostgres(database)
	relay := worker.NewOutboxRelay(outboxRepo, producer, worker.Config{
		WorkerCount:  cfg.OutboxWorkerCount,
		BatchSize:    cfg.OutboxBatchSize,
		PollInterval: cfg.OutboxPollInterval,
		MaxAttempts:  cfg.OutboxMaxAttempts,
		BaseBackoff:  cfg.OutboxBaseBackoff,
		MaxBackoff:   cfg.OutboxMaxBackoff,
	}, zapLog)

	var wg sync.WaitGroup
	wg.Go(func() {
		relay.Run(ctx)
	})

	repo := repository.NewRepoPostgres(database)
	svc := service.NewService(repo, rateProvider)
	h := handler.NewHandler(svc, zapLog)
	r := handler.NewRouter(h)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		zapLog.Info("server listening", zap.String("port", cfg.ServerPort))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zapLog.Fatal("server failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	zapLog.Info("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		zapLog.Error("http shutdown", zap.Error(err))
	}

	wg.Wait()
}
