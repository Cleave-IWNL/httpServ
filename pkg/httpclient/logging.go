package httpclient

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type LoggingClient struct {
	delegate HTTPClient
	logger   *zap.Logger
}

func NewLoggingClient(delegate HTTPClient, logger *zap.Logger) *LoggingClient {
	return &LoggingClient{
		delegate: delegate,
		logger:   logger,
	}
}

func (c *LoggingClient) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := c.delegate.Do(req)

	fields := []zap.Field{
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Duration("duration", time.Since(start)),
	}

	if err != nil {
		c.logger.Error("http request failed",
			append(fields, zap.Error(err))...)
		return resp, err
	}

	c.logger.Info("http request done",
		append(fields, zap.Int("status", resp.StatusCode))...)

	return resp, nil
}
