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

func (c *Client) Do(ctx context.Context, method, path string, headers map[string]string, payload, response any) error {
	var reqBody io.Reader = nil
	if payload != nil {
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		reqBody = strings.NewReader(string(jsonBytes))
	}

	req, err := http.NewRequestWithContext(ctx, method, c.config.BaseURL+path, reqBody)
	if err != nil {
		return err
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.config.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d with body %s", resp.StatusCode, string(body))
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(&response)
}

func NewClient(config Config) *Client {
	if config.Client == nil {
		config.Client = http.DefaultClient
	}

	return &Client{config: config}
}
