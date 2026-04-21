package httpclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type RetryConfig struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	ShouldRetry  func(*http.Response, error) bool
}

type RetryClient struct {
	delegate HTTPClient
	config   RetryConfig
	logger   *zap.Logger
}

func NewRetryClient(delegate HTTPClient, config RetryConfig, logger *zap.Logger) *RetryClient {
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	if config.InitialDelay == 0 {
		config.InitialDelay = 100 * time.Millisecond
	}

	if config.MaxDelay == 0 {
		config.MaxDelay = 5 * time.Second
	}

	if config.ShouldRetry == nil {
		config.ShouldRetry = defaultShouldRetry
	}

	return &RetryClient{
		delegate: delegate,
		config:   config,
		logger:   logger,
	}
}

func defaultShouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		return !errors.Is(err, context.Canceled) &&
			!errors.Is(err, context.DeadlineExceeded)
	}
	if resp == nil {
		return true
	}
	return resp.StatusCode >= 500 || resp.StatusCode == http.StatusTooManyRequests
}

func (c *RetryClient) Do(req *http.Request) (*http.Response, error) {
	var (
		resp  *http.Response
		err   error
		delay = c.config.InitialDelay
	)

	ctx := req.Context()

	for attempt := 0; attempt < c.config.MaxRetries; attempt++ {
		if req.GetBody != nil {
			body, bodyErr := req.GetBody()
			if bodyErr != nil {
				return nil, fmt.Errorf("get request body: %w", bodyErr)
			}
			req.Body = body
		}

		r := req.Clone(ctx)
		resp, err = c.delegate.Do(r)

		if !c.config.ShouldRetry(resp, err) || attempt == c.config.MaxRetries-1 {
			return resp, err
		}

		if resp != nil && resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}

		backoff := jitter(delay)
		c.logger.Warn("retry attempt",
			zap.Int("attempt", attempt+1),
			zap.Int("max", c.config.MaxRetries),
			zap.Duration("backoff", backoff),
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()),
			zap.Error(err),
		)

		if sleepErr := sleepCtx(ctx, backoff); sleepErr != nil {
			return resp, fmt.Errorf("retry sleep canceled: %w", sleepErr)
		}

		delay *= 2
		if delay > c.config.MaxDelay {
			delay = c.config.MaxDelay
		}
	}

	return resp, err
}

func jitter(d time.Duration) time.Duration {
	if d <= 0 {
		return 100 * time.Millisecond
	}
	return d/2 + time.Duration(rand.Int63n(int64(d)))
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
