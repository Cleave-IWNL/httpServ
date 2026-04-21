package httpclient

import (
	"errors"
	"fmt"
)

var (
	ErrUpstream = errors.New("httpclient: upstream error")
	ErrClient   = errors.New("httpclient: client error")
)

type HTTPError struct {
	StatusCode int
	URL        string
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("httpclient: %s returned %d: %s", e.URL, e.StatusCode, e.Body)
}

func (e *HTTPError) Unwrap() error {
	switch {
	case e.StatusCode >= 500:
		return ErrUpstream
	case e.StatusCode >= 400:
		return ErrClient
	default:
		return nil
	}
}
