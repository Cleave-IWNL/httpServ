package exchangerate

import (
	"context"
	"encoding/json"
	"fmt"
	"httpServ/pkg/httpclient"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	http    httpclient.HTTPClient
	baseURL string
	apiKey  string
}

func New(http httpclient.HTTPClient, baseUrl, apiKey string) *Client {
	return &Client{
		http:    http,
		baseURL: baseUrl,
		apiKey:  apiKey,
	}
}

func (c *Client) GetRate(ctx context.Context, from, to string) (Rate, error) {
	url := fmt.Sprintf("%s/v6/%s/pair/%s/%s", strings.TrimRight(c.baseURL, "/"),
		c.apiKey, strings.ToUpper(from), strings.ToUpper(to))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return Rate{}, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.http.Do(req)

	if err != nil {
		return Rate{}, fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return Rate{}, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return Rate{}, &httpclient.HTTPError{
			StatusCode: resp.StatusCode,
			URL:        url,
			Body:       string(body),
		}
	}

	var dto pairResponse
	if err := json.Unmarshal(body, &dto); err != nil {
		return Rate{}, fmt.Errorf("decode response: %w", err)
	}

	if dto.Result != "success" {
		return Rate{}, mapAPIError(dto.ErrorType)
	}

	return Rate{
		From:      dto.BaseCode,
		To:        dto.TargetCode,
		Value:     dto.ConversionRate,
		FetchedAt: dto.TimeLastUpdateUnix,
	}, nil
}
