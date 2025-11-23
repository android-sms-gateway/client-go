package smsgateway

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/android-sms-gateway/client-go/rest"
)

const BaseURL = "https://api.sms-gate.app/3rdparty/v1"
const settingsPath = "/settings"

// BASE_URL is deprecated, use BaseURL instead
//
//nolint:revive,staticcheck // backward compatibility
const BASE_URL = BaseURL

type Config struct {
	Client   *http.Client // Optional HTTP Client, defaults to `http.DefaultClient`
	BaseURL  string       // Optional base URL, defaults to `https://api.sms-gate.app/3rdparty/v1`
	User     string       // Required username
	Password string       // Required password
}

type Client struct {
	*rest.Client

	headers map[string]string
}

// NewClient creates a new instance of the API Client.
func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = BaseURL
	}

	return &Client{
		Client: rest.NewClient(rest.Config{
			Client:  config.Client,
			BaseURL: config.BaseURL,
		}),
		headers: map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(config.User+":"+config.Password)),
		},
	}
}

// Send enqueues a message for sending.
func (c *Client) Send(ctx context.Context, message Message, options ...SendOption) (MessageState, error) {
	opts := new(SendOptions).Apply(options...)
	path := "/messages?" + opts.ToURLValues().Encode()
	resp := new(MessageState)

	if err := c.Do(ctx, http.MethodPost, path, c.headers, &message, resp); err != nil {
		return *resp, fmt.Errorf("failed to send message: %w", err)
	}

	return *resp, nil
}

// GetState returns message state by ID.
func (c *Client) GetState(ctx context.Context, messageID string) (MessageState, error) {
	path := fmt.Sprintf("/messages/%s", url.PathEscape(messageID))
	resp := new(MessageState)

	if err := c.Do(ctx, http.MethodGet, path, c.headers, nil, resp); err != nil {
		return *resp, fmt.Errorf("failed to get message state: %w", err)
	}

	return *resp, nil
}

// ListDevices returns registered devices.
func (c *Client) ListDevices(ctx context.Context) ([]Device, error) {
	path := "/devices"
	var devices []Device

	if err := c.Do(ctx, http.MethodGet, path, c.headers, nil, &devices); err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}

	return devices, nil
}

// DeleteDevice removes a device by ID.
func (c *Client) DeleteDevice(ctx context.Context, id string) error {
	path := fmt.Sprintf("/devices/%s", url.PathEscape(id))

	if err := c.Do(ctx, http.MethodDelete, path, c.headers, nil, nil); err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}

	return nil
}

// CheckHealth returns service health status.
func (c *Client) CheckHealth(ctx context.Context) (HealthResponse, error) {
	path := "/health"
	resp := new(HealthResponse)

	if err := c.Do(ctx, http.MethodGet, path, c.headers, nil, resp); err != nil {
		return *resp, fmt.Errorf("failed to check health: %w", err)
	}

	return *resp, nil
}

// ExportInbox exports messages via webhooks.
func (c *Client) ExportInbox(ctx context.Context, req MessagesExportRequest) error {
	path := "/inbox/export"

	if err := c.Do(ctx, http.MethodPost, path, c.headers, &req, nil); err != nil {
		return fmt.Errorf("failed to export inbox: %w", err)
	}

	return nil
}

// GetLogs retrieves log entries.
func (c *Client) GetLogs(ctx context.Context, from, to time.Time) ([]LogEntry, error) {
	query := url.Values{}
	query.Set("from", from.Format(time.RFC3339))
	query.Set("to", to.Format(time.RFC3339))
	path := fmt.Sprintf("/logs?%s", query.Encode())
	var logs []LogEntry

	if err := c.Do(ctx, http.MethodGet, path, c.headers, nil, &logs); err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}

	return logs, nil
}

// GetSettings returns current settings.
func (c *Client) GetSettings(ctx context.Context) (DeviceSettings, error) {
	path := settingsPath
	resp := new(DeviceSettings)

	if err := c.Do(ctx, http.MethodGet, path, c.headers, nil, resp); err != nil {
		return *resp, fmt.Errorf("failed to get settings: %w", err)
	}

	return *resp, nil
}

// UpdateSettings partially updates settings.
func (c *Client) UpdateSettings(ctx context.Context, settings DeviceSettings) (DeviceSettings, error) {
	path := settingsPath
	resp := new(DeviceSettings)

	if err := c.Do(ctx, http.MethodPatch, path, c.headers, &settings, resp); err != nil {
		return *resp, fmt.Errorf("failed to update settings: %w", err)
	}

	return *resp, nil
}

// ReplaceSettings replaces all settings.
func (c *Client) ReplaceSettings(ctx context.Context, settings DeviceSettings) (DeviceSettings, error) {
	path := settingsPath
	resp := new(DeviceSettings)

	if err := c.Do(ctx, http.MethodPut, path, c.headers, &settings, resp); err != nil {
		return *resp, fmt.Errorf("failed to replace settings: %w", err)
	}

	return *resp, nil
}

// ListWebhooks returns registered webhooks
// Returns a slice of Webhook objects or an error if the request fails.
func (c *Client) ListWebhooks(ctx context.Context) ([]Webhook, error) {
	path := "/webhooks"
	resp := []Webhook{}

	if err := c.Do(ctx, http.MethodGet, path, c.headers, nil, &resp); err != nil {
		return resp, fmt.Errorf("failed to list webhooks: %w", err)
	}

	return resp, nil
}

// RegisterWebhook registers or replaces a webhook
// Returns the registered webhook with server-assigned fields or an error if the request fails.
func (c *Client) RegisterWebhook(ctx context.Context, webhook Webhook) (Webhook, error) {
	path := "/webhooks"
	resp := new(Webhook)

	if err := c.Do(ctx, http.MethodPost, path, c.headers, &webhook, resp); err != nil {
		return *resp, fmt.Errorf("failed to register webhook: %w", err)
	}

	return *resp, nil
}

// DeleteWebhook removes a webhook by ID
// Returns an error if the deletion fails.
func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) error {
	path := fmt.Sprintf("/webhooks/%s", url.PathEscape(webhookID))

	if err := c.Do(ctx, http.MethodDelete, path, c.headers, nil, nil); err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	return nil
}
