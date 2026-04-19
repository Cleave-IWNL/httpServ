package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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
}

func NewRetryClient(delegate HTTPClient, config RetryConfig) *RetryClient {
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
		config.ShouldRetry = func(resp *http.Response, err error) bool {
			if err != nil {
				return true
			}

			return resp.StatusCode >= 500
		}
	}

	return &RetryClient{
		delegate: delegate,
		config:   config,
	}
}

func (c *RetryClient) Do(req *http.Request) (*http.Response, error) {
	var (
		resp  *http.Response
		err   error
		delay = c.config.InitialDelay
	)

	for i := 0; i < c.config.MaxRetries; i++ {
		if req.Body != nil {
			if seeker, ok := req.Body.(io.Seeker); ok {
				_, seekErr := seeker.Seek(0, io.SeekStart)

				if seekErr != nil {
					return nil, fmt.Errorf("failed to seek request body: %w", seekErr)
				}
			} else {
				bodyBytes, readErr := io.ReadAll(req.Body)

				if readErr != nil {
					return nil, fmt.Errorf("failed to read request body for retry: %w", readErr)
				}

				req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}

		resp, err = c.delegate.Do(req)

		if c.config.ShouldRetry(resp, err) {
			log.Printf("Request failed (attempt %d/%d), retrying in %v. Error: %v", i+1, c.config.MaxRetries, delay, err)
			time.Sleep(delay)
			delay = min(time.Duration(float64(delay)*2), c.config.MaxDelay)
			continue
		}
		return resp, err
	}

	return resp, err
}
