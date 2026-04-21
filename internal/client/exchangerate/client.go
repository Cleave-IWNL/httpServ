package exchangerate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"httpServ/internal/model"
	"httpServ/pkg/httpclient"
)

const (
	maxResponseBodyBytes = 64 * 1024
	maxErrorBodyBytes    = 4 * 1024
)

type Client struct {
	http    httpclient.HTTPClient
	baseURL string
	apiKey  string
}

func New(httpClient httpclient.HTTPClient, baseURL, apiKey string) *Client {
	return &Client{
		http:    httpClient,
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
	}
}

func (c *Client) GetRate(ctx context.Context, from, to string) (model.Rate, error) {
	url := fmt.Sprintf("%s/v6/%s/pair/%s/%s",
		c.baseURL, c.apiKey,
		strings.ToUpper(from), strings.ToUpper(to))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.Rate{}, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return model.Rate{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodyBytes))
	if err != nil {
		return model.Rate{}, fmt.Errorf("read body: %w", err)
	}

	var dto model.ExchangeRatePair
	parseErr := json.Unmarshal(body, &dto)

	if parseErr == nil && dto.Result == "error" {
		return model.Rate{}, mapAPIError(dto.ErrorType)
	}

	if resp.StatusCode != http.StatusOK {
		snippet := body
		if len(snippet) > maxErrorBodyBytes {
			snippet = snippet[:maxErrorBodyBytes]
		}
		return model.Rate{}, &httpclient.HTTPError{
			StatusCode: resp.StatusCode,
			URL:        url,
			Body:       string(snippet),
		}
	}

	if parseErr != nil {
		return model.Rate{}, fmt.Errorf("decode response: %w", parseErr)
	}

	if dto.Result != "success" {
		return model.Rate{}, mapAPIError(dto.ErrorType)
	}

	return model.Rate{
		From:      dto.BaseCode,
		To:        dto.TargetCode,
		Value:     dto.ConversionRate,
		FetchedAt: dto.TimeLastUpdateUnix,
	}, nil
}
