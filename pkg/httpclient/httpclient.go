package httpclient

import (
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type defaultClient struct {
	client *http.Client
}

func NewDefaultClient(timeout time.Duration) HTTPClient {
	return &defaultClient{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
	}
}

func (c *defaultClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
