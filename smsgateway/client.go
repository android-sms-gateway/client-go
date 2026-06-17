package smsgateway

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/android-sms-gateway/client-go/rest"
)

const BaseURL = "https://api.sms-gate.app/3rdparty/v1"
const settingsPath = "/settings"

// Deprecated: BASE_URL is kept for backward compatibility. Use BaseURL instead.
//
//nolint:revive,staticcheck // backward compatibility
const BASE_URL = BaseURL

type Client struct {
	*rest.Client

	headers map[string]string
}

// NewClient creates a new instance of the API Client.
func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = BaseURL
	}

	headers := make(map[string]string, 1)
	if config.Token != "" {
		headers["Authorization"] = "Bearer " + config.Token
	} else {
		headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(config.User+":"+config.Password))
	}

	return &Client{
		Client: rest.NewClient(rest.Config{
			Client:  config.Client,
			BaseURL: config.BaseURL,
		}),
		headers: headers,
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

// CancelMessage cancels a pending message by ID.
func (c *Client) CancelMessage(ctx context.Context, messageID string) error {
	path := fmt.Sprintf("/messages/%s", url.PathEscape(messageID))

	if err := c.Do(ctx, http.MethodDelete, path, c.headers, nil, nil); err != nil {
		return fmt.Errorf("failed to cancel message: %w", err)
	}

	return nil
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
//
// Deprecated: use RefreshInbox instead.
func (c *Client) ExportInbox(ctx context.Context, req MessagesExportRequest) error {
	path := "/inbox/export"

	if err := c.Do(ctx, http.MethodPost, path, c.headers, &req, nil); err != nil {
		return fmt.Errorf("failed to export inbox: %w", err)
	}

	return nil
}

// ListInboxMessages retrieves incoming messages with filtering and pagination.
// Returns the messages, total count (from X-Total-Count header), and error.
func (c *Client) ListInboxMessages(ctx context.Context, opts ListInboxOptions) ([]IncomingMessage, int, error) {
	path := "/inbox?" + opts.ToURLValues().Encode()
	var msgs []IncomingMessage

	hdr, err := c.DoWithResponseHeaders(ctx, http.MethodGet, path, c.headers, nil, &msgs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list inbox messages: %w", err)
	}

	total := 0
	if v := hdr.Get("X-Total-Count"); v != "" {
		total, err = strconv.Atoi(v)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse X-Total-Count header: %w", err)
		}
	}

	return msgs, total, nil
}

// RefreshInbox requests an inbox messages refresh.
func (c *Client) RefreshInbox(ctx context.Context, req InboxRefreshRequest) error {
	path := "/inbox/refresh"

	if err := c.Do(ctx, http.MethodPost, path, c.headers, &req, nil); err != nil {
		return fmt.Errorf("failed to refresh inbox: %w", err)
	}

	return nil
}

// ListMessages retrieves messages with filtering and pagination.
// Returns the messages, total count (from X-Total-Count header), and error.
func (c *Client) ListMessages(ctx context.Context, opts ListMessagesOptions) ([]MessageState, int, error) {
	path := "/messages?" + opts.ToURLValues().Encode()
	var msgs []MessageState

	hdr, err := c.DoWithResponseHeaders(ctx, http.MethodGet, path, c.headers, nil, &msgs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list messages: %w", err)
	}

	total := 0
	if v := hdr.Get("X-Total-Count"); v != "" {
		total, err = strconv.Atoi(v)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse X-Total-Count header: %w", err)
		}
	}

	return msgs, total, nil
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

// GenerateToken generates a new access token with specified scopes and ttl.
// Returns the generated token details or an error if the request fails.
func (c *Client) GenerateToken(ctx context.Context, req TokenRequest) (TokenResponse, error) {
	path := "/auth/token"
	resp := new(TokenResponse)

	if err := c.Do(ctx, http.MethodPost, path, c.headers, &req, resp); err != nil {
		return *resp, fmt.Errorf("failed to generate token: %w", err)
	}

	return *resp, nil
}

// RefreshToken exchanges a refresh token for a new token pair.
// Returns the refreshed token details or an error if the request fails.
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (TokenResponse, error) {
	path := "/auth/token/refresh"
	resp := new(TokenResponse)
	headers := map[string]string{
		"Authorization": "Bearer " + refreshToken,
	}

	if err := c.Do(ctx, http.MethodPost, path, headers, nil, resp); err != nil {
		return *resp, fmt.Errorf("failed to refresh token: %w", err)
	}

	return *resp, nil
}

// RevokeToken revokes an access token with the specified jti (token ID).
// Returns an error if the revocation fails.
func (c *Client) RevokeToken(ctx context.Context, jti string) error {
	path := fmt.Sprintf("/auth/token/%s", url.PathEscape(jti))

	if err := c.Do(ctx, http.MethodDelete, path, c.headers, nil, nil); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}
