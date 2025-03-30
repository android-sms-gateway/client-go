package smsgateway_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

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

		if string(req) == `{"message":"","phoneNumbers":null}` {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(req)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := smsgateway.NewClient(smsgateway.Config{
		BaseURL: server.URL,
	})

	type args struct {
		ctx     context.Context
		message smsgateway.Message
	}
	tests := []struct {
		name    string
		c       *smsgateway.Client
		args    args
		want    smsgateway.MessageState
		wantErr bool
	}{
		{
			name: "Success",
			c:    client,
			args: args{
				ctx: context.TODO(),
				message: smsgateway.Message{
					Message:      "Hello, world!",
					PhoneNumbers: []string{"+1234567890"},
				},
			},
			want:    smsgateway.MessageState{},
			wantErr: false,
		},
		{
			name: "Bad Request",
			c:    client,
			args: args{
				ctx:     context.TODO(),
				message: smsgateway.Message{},
			},
			want:    smsgateway.MessageState{},
			wantErr: true,
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
	t.Run("Successful request", func(t *testing.T) {
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

		client := smsgateway.NewClient(smsgateway.Config{
			BaseURL:  server.URL,
			Client:   nil,
			User:     "",
			Password: "",
		})

		state, err := client.GetState(context.Background(), "123")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if state.ID != "123" {
			t.Errorf("Expected ID 123, got %s", state.ID)
		}
		if state.State != smsgateway.ProcessingStatePending {
			t.Errorf("Expected state Pending, got %s", state.State)
		}
	})

	// Test case 2: Error response
	t.Run("Error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := smsgateway.NewClient(smsgateway.Config{
			BaseURL:  server.URL,
			Client:   nil,
			User:     "",
			Password: "",
		})

		_, err := client.GetState(context.Background(), "123")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	// Test case 3: Authorization header present
	t.Run("Authorization header present", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Basic dXNlcjpwYXNzd29yZA==" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id": "123", "state": "Pending"}`))
		}))
		defer server.Close()

		client := smsgateway.NewClient(smsgateway.Config{
			BaseURL:  server.URL,
			Client:   nil,
			User:     "user",
			Password: "password",
		})

		_, err := client.GetState(context.Background(), "123")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
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

	client := smsgateway.NewClient(smsgateway.Config{
		BaseURL:  server.URL,
		Client:   nil,
		User:     "",
		Password: "",
	})

	tests := []struct {
		name    string
		c       *smsgateway.Client
		want    []smsgateway.Webhook
		wantErr bool
	}{
		{
			name: "Success",
			c:    client,
			want: []smsgateway.Webhook{
				{
					ID:    "123",
					URL:   "https://example.com",
					Event: smsgateway.WebhookEventSmsDelivered,
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

	client := smsgateway.NewClient(smsgateway.Config{
		BaseURL:  server.URL,
		Client:   nil,
		User:     "",
		Password: "",
	})

	type args struct {
		webhook smsgateway.Webhook
	}
	tests := []struct {
		name    string
		c       *smsgateway.Client
		args    args
		want    smsgateway.Webhook
		wantErr bool
	}{
		{
			name: "Success",
			c:    client,
			args: args{
				webhook: smsgateway.Webhook{
					ID:    "",
					URL:   "https://example.com",
					Event: smsgateway.WebhookEventSmsDelivered,
				},
			},
			want: smsgateway.Webhook{
				ID:    "123",
				URL:   "https://example.com",
				Event: smsgateway.WebhookEventSmsDelivered,
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

	client := smsgateway.NewClient(smsgateway.Config{
		BaseURL:  server.URL,
		Client:   nil,
		User:     "",
		Password: "",
	})

	type args struct {
		webhookID string
	}
	tests := []struct {
		name    string
		c       *smsgateway.Client
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
