package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Config struct {
	Client  *http.Client // Optional HTTP Client, defaults to `http.DefaultClient`
	BaseURL string       // Optional base URL
}

type Client struct {
	config Config
}

func NewClient(config Config) *Client {
	if config.Client == nil {
		config.Client = http.DefaultClient
	}

	return &Client{config: config}
}

func (c *Client) Do(ctx context.Context, method, path string, headers map[string]string, payload, response any) error {
	_, err := c.do(ctx, method, path, headers, payload, response)
	return err
}

// DoWithResponseHeaders is like Do but also returns the response headers.
func (c *Client) DoWithResponseHeaders(
	ctx context.Context,
	method, path string,
	headers map[string]string,
	payload, response any,
) (http.Header, error) {
	return c.do(ctx, method, path, headers, payload, response)
}

func (c *Client) do(
	ctx context.Context,
	method, path string,
	headers map[string]string,
	payload, response any,
) (http.Header, error) {
	var reqBody io.Reader
	if payload != nil {
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		reqBody = strings.NewReader(string(jsonBytes))
	}

	req, err := http.NewRequestWithContext(ctx, method, c.config.BaseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.config.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)

		return resp.Header, c.formatError(resp.StatusCode, body)
	}

	if resp.StatusCode == http.StatusNoContent {
		return resp.Header, nil
	}

	if response != nil {
		if decErr := json.NewDecoder(resp.Body).Decode(response); decErr != nil {
			return nil, fmt.Errorf("failed to decode response: %w", decErr)
		}
	}

	return resp.Header, nil
}

func (c *Client) formatError(statusCode int, body []byte) error {
	switch statusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", ErrBadRequest, string(body))
	case http.StatusConflict:
		return fmt.Errorf("%w: %s", ErrConflict, string(body))
	}

	if statusCode >= http.StatusInternalServerError {
		return fmt.Errorf("%w: unexpected status code %d with body %s", ErrServer, statusCode, string(body))
	}

	// All other client errors (400-499)
	return fmt.Errorf("%w: unexpected status code %d with body %s", ErrClient, statusCode, string(body))
}
