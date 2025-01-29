package smsgateway

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	// Test case 1: Test with default client and base URL
	config := Config{}
	client := NewClient(config)
	if client.config.Client != http.DefaultClient {
		t.Errorf("Expected default client, got %v", client.config.Client)
	}
	if client.config.BaseURL != BASE_URL {
		t.Errorf("Expected default base URL, got %s", client.config.BaseURL)
	}

	// Test case 2: Test with custom client and base URL
	customClient := &http.Client{}
	customBaseURL := "https://example.com"
	config = Config{Client: customClient, BaseURL: customBaseURL}
	client = NewClient(config)
	if client.config.Client != customClient {
		t.Errorf("Expected custom client, got %v", client.config.Client)
	}
	if client.config.BaseURL != customBaseURL {
		t.Errorf("Expected custom base URL, got %s", client.config.BaseURL)
	}
}

func TestClient_Send(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/message" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		req, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		if string(req) != `{"message":"","phoneNumbers":null}` {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(req)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
	})

	type args struct {
		ctx     context.Context
		message Message
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		want    MessageState
		wantErr bool
	}{
		{
			name: "Success",
			c:    client,
			args: args{
				ctx:     context.TODO(),
				message: Message{},
			},
			want:    MessageState{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Send(tt.args.ctx, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Send() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetState(t *testing.T) {
	// Test case 1: Successful request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/message/123" {
			t.Errorf("Expected path /message/123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id": "123", "state": "Pending"}`))
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
	})

	state, err := client.GetState(context.Background(), "123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if state.ID != "123" {
		t.Errorf("Expected ID 123, got %s", state.ID)
	}
	if state.State != ProcessingStatePending {
		t.Errorf("Expected state Pending, got %s", state.State)
	}

	// Test case 2: Error response
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client = NewClient(Config{
		BaseURL: server.URL,
	})

	_, err = client.GetState(context.Background(), "123")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Test case 3: Invalid message ID
	_, err = client.GetState(context.Background(), "invalid")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestClient_ListWebhooks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/webhooks" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":"123","url":"https://example.com","event":"sms:delivered"}]`))
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
	})

	tests := []struct {
		name    string
		c       *Client
		want    []Webhook
		wantErr bool
	}{
		{
			name: "Success",
			c:    client,
			want: []Webhook{
				{
					ID:    "123",
					URL:   "https://example.com",
					Event: WebhookEventSmsDelivered,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.ListWebhooks(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListWebhooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListWebhooks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_RegisterWebhook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/webhooks" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		if string(body) != `{"url":"https://example.com","event":"sms:delivered"}` {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"123","url":"https://example.com","event":"sms:delivered"}`))
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
	})

	type args struct {
		webhook Webhook
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		want    Webhook
		wantErr bool
	}{
		{
			name: "Success",
			c:    client,
			args: args{
				webhook: Webhook{
					URL:   "https://example.com",
					Event: WebhookEventSmsDelivered,
				},
			},
			want: Webhook{
				ID:    "123",
				URL:   "https://example.com",
				Event: WebhookEventSmsDelivered,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.RegisterWebhook(context.Background(), tt.args.webhook)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RegisterWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.RegisterWebhook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DeleteWebhook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/webhooks/123" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
	})

	type args struct {
		webhookID string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			c:    client,
			args: args{
				webhookID: "123",
			},
			wantErr: false,
		},
		{
			name: "Not Found",
			c:    client,
			args: args{
				webhookID: "456",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.DeleteWebhook(context.Background(), tt.args.webhookID); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
