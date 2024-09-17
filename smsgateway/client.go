package smsgateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	return json.NewDecoder(resp.Body).Decode(&response)
}
