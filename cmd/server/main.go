package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"httpServ/internal/client/exchangerate"
	"httpServ/internal/handler"
	"httpServ/internal/repository"
	"httpServ/internal/service"
	"httpServ/pkg/config"
	"httpServ/pkg/db"
	"httpServ/pkg/httpclient"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	relay := worker.NewOutboxRelay(database, zapLog)
	go relay.Run(ctx)

	baseHTTP := httpclient.NewDefaultClient(5 * time.Second)
	retryHTTP := httpclient.NewRetryClient(baseHTTP, httpclient.RetryConfig{
		MaxRetries: 3,
	})
	loggedHTTP := httpclient.NewLoggingClient(retryHTTP, zapLog)

	rateProvider := exchangerate.New(loggedHTTP, cfg.ExchangeAPIURL, cfg.ExchangeAPIKey)

	repo := repository.NewRepoPostgres(database)
	svc := service.NewService(repo, rateProvider)
	h := handler.NewHandler(svc, zapLog)
	r := handler.NewRouter(h)

	if err = http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		zapLog.Fatal("server failed", zap.Error(err))
	}
}
