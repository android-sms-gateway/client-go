package smsgateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const BASE_URL = "https://api.sms-gate.app/3rdparty/v1"

type Config struct {
	Client   *http.Client // Optional HTTP Client, defaults to `http.DefaultClient`
	BaseURL  string       // Optional base URL, defaults to `https://api.sms-gate.app/3rdparty/v1`
	User     string       // Required username
	Password string       // Required password
}

type Client struct {
	config Config
}

// NewClient creates a new instance of the API Client.
func NewClient(config Config) *Client {
	if config.Client == nil {
		config.Client = http.DefaultClient
	}
	if config.BaseURL == "" {
		config.BaseURL = BASE_URL
	}

	return &Client{config: config}
}

// Sends an SMS message.
func (c *Client) Send(ctx context.Context, message Message) (MessageState, error) {
	path := "/message"
	resp := MessageState{}

	return resp, c.doRequest(ctx, http.MethodPost, path, map[string]string{}, &message, &resp)
}

// Gets the state of an SMS message by ID.
func (c *Client) GetState(ctx context.Context, messageID string) (MessageState, error) {
	path := fmt.Sprintf("/message/%s", messageID)
	resp := MessageState{}

	return resp, c.doRequest(ctx, http.MethodGet, path, map[string]string{}, nil, &resp)
}

// ListWebhooks retrieves all registered webhooks.
// Returns a slice of Webhook objects or an error if the request fails.
func (c *Client) ListWebhooks(ctx context.Context) ([]Webhook, error) {
	path := "/webhooks"
	resp := []Webhook{}

	return resp, c.doRequest(ctx, http.MethodGet, path, map[string]string{}, nil, &resp)
}

// RegisterWebhook registers a new webhook.
// Returns the registered webhook with server-assigned fields or an error if the request fails.
func (c *Client) RegisterWebhook(ctx context.Context, webhook Webhook) (Webhook, error) {
	path := "/webhooks"
	resp := Webhook{}

	return resp, c.doRequest(ctx, http.MethodPost, path, map[string]string{}, &webhook, &resp)
}

// DeleteWebhook removes a webhook with the specified ID.
// Returns an error if the deletion fails.
func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) error {
	path := fmt.Sprintf("/webhooks/%s", url.PathEscape(webhookID))

	return c.doRequest(ctx, http.MethodDelete, path, map[string]string{}, nil, nil)
}

func (c *Client) doRequest(ctx context.Context, method, path string, headers map[string]string, payload, response any) error {
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

	req.SetBasicAuth(c.config.User, c.config.Password)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := c.config.Client.Do(req)
	if err != nil {
		return err
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
