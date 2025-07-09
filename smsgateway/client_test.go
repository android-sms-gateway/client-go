package smsgateway_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

const (
	username            = "username"
	password            = "password"
	authorizationHeader = "Basic dXNlcm5hbWU6cGFzc3dvcmQ="
)

type mockServerExpectedInput struct {
	method      string
	path        string
	query       string
	contentType string
	body        string
}

type mockServerOutput struct {
	code int
	body string
}

func newMockServer(input mockServerExpectedInput, output mockServerOutput) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != input.method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != input.path {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if input.query != "" && r.URL.RawQuery != input.query {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.Header.Get("Authorization") != authorizationHeader {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if input.contentType != "" && r.Header.Get("Content-Type") != input.contentType {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		req, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		if input.body != "" && string(req) != input.body {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(output.code)
		_, _ = w.Write([]byte(output.body))
	}))
}

func newClient(baseURL string) *smsgateway.Client {
	return smsgateway.NewClient(smsgateway.Config{
		BaseURL:  baseURL,
		User:     username,
		Password: password,
	})
}

func TestClient_Send(t *testing.T) {
	type args struct {
		ctx     context.Context
		message smsgateway.Message
		options []smsgateway.SendOption
	}
	tests := []struct {
		name    string
		args    args
		want    smsgateway.MessageState
		wantErr bool
		query   string
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				message: smsgateway.Message{
					TextMessage: &smsgateway.TextMessage{
						Text: "Hello World!",
					},
					PhoneNumbers: []string{"+1234567890"},
				},
			},
			want:    smsgateway.MessageState{},
			wantErr: false,
		},
		{
			name: "Bad Request",
			args: args{
				ctx:     context.Background(),
				message: smsgateway.Message{},
			},
			want:    smsgateway.MessageState{},
			wantErr: true,
		},
		{
			name: "WithSkipPhoneValidation=true",
			args: args{
				ctx: context.Background(),
				message: smsgateway.Message{
					TextMessage: &smsgateway.TextMessage{
						Text: "Hello World!",
					},
					PhoneNumbers: []string{"+1234567890"},
				},
				options: []smsgateway.SendOption{
					smsgateway.WithSkipPhoneValidation(true),
				},
			},
			want:    smsgateway.MessageState{},
			wantErr: false,
			query:   "skipPhoneValidation=true",
		},
		{
			name: "WithSkipPhoneValidation=false",
			args: args{
				ctx: context.Background(),
				message: smsgateway.Message{
					TextMessage: &smsgateway.TextMessage{
						Text: "Hello World!",
					},
					PhoneNumbers: []string{"+1234567890"},
				},
				options: []smsgateway.SendOption{
					smsgateway.WithSkipPhoneValidation(false),
				},
			},
			want:    smsgateway.MessageState{},
			wantErr: false,
			query:   "skipPhoneValidation=false",
		},
		{
			name: "WithDeviceActiveWithin",
			args: args{
				ctx: context.Background(),
				message: smsgateway.Message{
					TextMessage: &smsgateway.TextMessage{
						Text: "Hello World!",
					},
					PhoneNumbers: []string{"+1234567890"},
				},
				options: []smsgateway.SendOption{
					smsgateway.WithDeviceActiveWithin(24),
				},
			},
			want:    smsgateway.MessageState{},
			wantErr: false,
			query:   "deviceActiveWithin=24",
		},
		{
			name: "WithBothOptions",
			args: args{
				ctx: context.Background(),
				message: smsgateway.Message{
					TextMessage: &smsgateway.TextMessage{
						Text: "Hello World!",
					},
					PhoneNumbers: []string{"+1234567890"},
				},
				options: []smsgateway.SendOption{
					smsgateway.WithSkipPhoneValidation(true),
					smsgateway.WithDeviceActiveWithin(48),
				},
			},
			want:    smsgateway.MessageState{},
			wantErr: false,
			query:   "deviceActiveWithin=48&skipPhoneValidation=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock server for each test to validate query parameters
			server := newMockServer(mockServerExpectedInput{
				method:      http.MethodPost,
				path:        "/messages",
				contentType: "application/json",
				body:        `{"textMessage":{"text":"Hello World!"},"phoneNumbers":["+1234567890"]}`,
				query:       tt.query,
			}, mockServerOutput{
				code: http.StatusCreated,
				body: `{}`,
			})
			defer server.Close()

			client := newClient(server.URL)

			got, err := client.Send(tt.args.ctx, tt.args.message, tt.args.options...)
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
		server := newMockServer(mockServerExpectedInput{
			method: http.MethodGet,
			path:   "/messages/123",
		}, mockServerOutput{
			code: http.StatusOK,
			body: `{"id": "123", "state": "Pending"}`,
		},
		)
		defer server.Close()

		client := newClient(server.URL)

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
		server := newMockServer(mockServerExpectedInput{
			method: http.MethodGet,
			path:   "/messages/123",
		}, mockServerOutput{
			code: http.StatusInternalServerError,
			body: `{"error": "Internal server error"}`,
		},
		)
		defer server.Close()

		client := newClient(server.URL)

		_, err := client.GetState(context.Background(), "123")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestClient_ListWebhooks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := newMockServer(mockServerExpectedInput{
			method: http.MethodGet,
			path:   "/webhooks",
		}, mockServerOutput{
			code: http.StatusOK,
			body: `[{"id":"123","deviceId":null,"url":"https://example.com","event":"sms:delivered"}]`,
		},
		)
		defer server.Close()

		client := newClient(server.URL)

		res, err := client.ListWebhooks(context.Background())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := []smsgateway.Webhook{
			{
				ID:    "123",
				URL:   "https://example.com",
				Event: smsgateway.WebhookEventSmsDelivered,
			},
		}

		if !reflect.DeepEqual(res, expected) {
			t.Errorf("Expected %v, got %v", expected, res)
		}
	})

	t.Run("Error response", func(t *testing.T) {
		server := newMockServer(mockServerExpectedInput{
			method: http.MethodGet,
			path:   "/webhooks",
		}, mockServerOutput{
			code: http.StatusInternalServerError,
			body: `{"error": "Internal server error"}`,
		},
		)
		defer server.Close()

		client := newClient(server.URL)

		_, err := client.ListWebhooks(context.Background())
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestClient_RegisterWebhook(t *testing.T) {
	server := newMockServer(mockServerExpectedInput{
		method:      http.MethodPost,
		path:        "/webhooks",
		contentType: "application/json",
		body:        `{"deviceId":null,"url":"https://example.com","event":"sms:delivered"}`,
	}, mockServerOutput{
		code: http.StatusCreated,
		body: `{"id":"123","deviceId":null,"url":"https://example.com","event":"sms:delivered"}`,
	})
	defer server.Close()

	client := newClient(server.URL)

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
		{
			name: "Error response",
			c:    client,
			args: args{
				webhook: smsgateway.Webhook{},
			},
			want:    smsgateway.Webhook{},
			wantErr: true,
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
	server := newMockServer(mockServerExpectedInput{
		method: http.MethodDelete,
		path:   "/webhooks/123",
	}, mockServerOutput{
		code: http.StatusNoContent,
	})
	defer server.Close()

	client := newClient(server.URL)

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

func TestClient_ListDevices(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := newMockServer(mockServerExpectedInput{
			method: http.MethodGet,
			path:   "/devices",
		}, mockServerOutput{
			code: http.StatusOK,
			body: `[{"createdAt":"2025-02-14T07:00:52.245+03:00","id":"7u21IKyeZL0uyKot196nY","lastSeen":"2025-02-14T07:10:05.732+03:00","name":"Google/sdk_gphone64_arm64","updatedAt":"2025-02-14T07:10:05.739+03:00"},{"createdAt":"2024-02-03T16:50:20.766+03:00","id":"ALsLlZFmwkYxkSn_yovQT","lastSeen":"2025-02-12T12:22:47.673+03:00","name":"Google/sdk_gphone64_arm64","updatedAt":"2025-02-12T12:22:47.677+03:00"}]`,
		})
		defer server.Close()

		client := newClient(server.URL)

		devices, err := client.ListDevices(context.Background())
		if err != nil {
			t.Fatalf("ListDevices failed: %v", err)
		}
		if len(devices) != 2 {
			t.Errorf("Expected 2 devices, got %d", len(devices))
		}
	})

	t.Run("Error response", func(t *testing.T) {
		server := newMockServer(mockServerExpectedInput{
			method: http.MethodGet,
			path:   "/devices",
		}, mockServerOutput{
			code: http.StatusInternalServerError,
		})
		defer server.Close()

		client := newClient(server.URL)

		_, err := client.ListDevices(context.Background())
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}

func TestClient_DeleteDevice(t *testing.T) {
	server := newMockServer(mockServerExpectedInput{
		method: http.MethodDelete,
		path:   "/devices/123",
	}, mockServerOutput{
		code: http.StatusNoContent,
	})
	defer server.Close()

	client := newClient(server.URL)

	type args struct {
		deviceID string
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
				deviceID: "123",
			},
			wantErr: false,
		},
		{
			name: "Not Found",
			c:    client,
			args: args{
				deviceID: "456",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.DeleteDevice(context.Background(), tt.args.deviceID); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteDevice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CheckHealth(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := newMockServer(mockServerExpectedInput{
			method: http.MethodGet,
			path:   "/health",
		}, mockServerOutput{
			code: http.StatusOK,
			body: `{"checks":{"db:ping":{"description":"Failed sequential pings count","observedValue":0,"status":"pass"}},"releaseId":1117,"status":"pass","version":"v1.24.0"}`,
		})
		defer server.Close()

		client := newClient(server.URL)

		health, err := client.CheckHealth(context.Background())
		if err != nil {
			t.Fatalf("CheckHealth failed: %v", err)
		}
		if health.Status != "pass" {
			t.Errorf("Expected status 'pass', got '%s'", health.Status)
		}
	})

	t.Run("InternalError", func(t *testing.T) {
		server := newMockServer(mockServerExpectedInput{
			method: http.MethodGet,
			path:   "/health",
		}, mockServerOutput{
			code: http.StatusInternalServerError,
		})
		defer server.Close()

		client := newClient(server.URL)

		_, err := client.CheckHealth(context.Background())
		if err == nil {
			t.Fatal("Expected error for internal server error")
		}
	})
}

func TestClient_ExportInbox(t *testing.T) {
	server := newMockServer(mockServerExpectedInput{
		method: http.MethodPost,
		path:   "/inbox/export",
		body:   `{"deviceId":"dev1","since":"2024-01-01T00:00:00Z","until":"2024-01-02T00:00:00Z"}`,
	}, mockServerOutput{
		code: http.StatusNoContent,
	})
	defer server.Close()

	client := newClient(server.URL)

	tests := []struct {
		name    string
		args    smsgateway.MessagesExportRequest
		wantErr bool
	}{
		{
			name: "Success",
			args: smsgateway.MessagesExportRequest{
				DeviceID: "dev1",
				Since:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Until:    time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "Invalid request",
			args: smsgateway.MessagesExportRequest{
				DeviceID: "dev1",
				Since:    time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				Until:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := client.ExportInbox(context.Background(), tt.args); (err != nil) != tt.wantErr {
				t.Errorf("Client.ExportInbox() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetLogs(t *testing.T) {
	server := newMockServer(mockServerExpectedInput{
		method: http.MethodGet,
		path:   "/logs",
		query:  "from=2024-01-01T00%3A00%3A00Z&to=2024-01-02T00%3A00%3A00Z",
	}, mockServerOutput{
		code: http.StatusOK,
		body: `[{"id":1,"message":"Test log"}]`,
	})
	defer server.Close()

	client := newClient(server.URL)

	tests := []struct {
		name string
		args struct {
			from time.Time
			to   time.Time
		}
		want    []smsgateway.LogEntry
		wantErr bool
	}{
		{
			name: "Success",
			args: struct {
				from time.Time
				to   time.Time
			}{
				from: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			want: []smsgateway.LogEntry{
				{
					ID:      1,
					Message: "Test log",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid request",
			args: struct {
				from time.Time
				to   time.Time
			}{
				from: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetLogs(context.Background(), tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetLogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetLogs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetSettings(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := newMockServer(mockServerExpectedInput{
			method: http.MethodGet,
			path:   "/settings",
		}, mockServerOutput{
			code: http.StatusOK,
			body: `{"messages":{"limit_period":"PerDay","limit_value":100}}`,
		})
		defer server.Close()

		client := newClient(server.URL)

		settings, err := client.GetSettings(context.Background())
		if err != nil {
			t.Fatalf("GetSettings failed: %v", err)
		}
		if *settings.Messages.LimitPeriod != smsgateway.PerDay {
			t.Errorf("Expected limit period 'PerDay', got '%v'", *settings.Messages.LimitPeriod)
		}
	})

	t.Run("Error", func(t *testing.T) {
		server := newMockServer(mockServerExpectedInput{
			method: http.MethodGet,
			path:   "/settings",
		}, mockServerOutput{
			code: http.StatusInternalServerError,
		})
		defer server.Close()

		client := newClient(server.URL)

		_, err := client.GetSettings(context.Background())
		if err == nil {
			t.Fatal("Expected error for internal server error")
		}
	})
}

func TestClient_UpdateSettings(t *testing.T) {
	server := newMockServer(mockServerExpectedInput{
		method:      http.MethodPatch,
		path:        "/settings",
		contentType: "application/json",
		body:        `{"messages":{"limit_period":"PerHour","limit_value":50}}`,
	}, mockServerOutput{
		code: http.StatusOK,
		body: `{"messages":{"limit_period":"PerHour","limit_value":50}}`,
	})
	defer server.Close()

	client := newClient(server.URL)

	limitPeriod := smsgateway.PerHour
	limitValue := 50
	tests := []struct {
		name     string
		args     smsgateway.DeviceSettings
		expected smsgateway.DeviceSettings
		wantErr  bool
	}{
		{
			name: "Success",
			args: smsgateway.DeviceSettings{
				Messages: &smsgateway.SettingsMessages{
					LimitPeriod: &limitPeriod,
					LimitValue:  &limitValue,
				},
			},
			expected: smsgateway.DeviceSettings{
				Messages: &smsgateway.SettingsMessages{
					LimitPeriod: &limitPeriod,
					LimitValue:  &limitValue,
				},
			},
			wantErr: false,
		},
		{
			name:     "Error",
			args:     smsgateway.DeviceSettings{},
			expected: smsgateway.DeviceSettings{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.UpdateSettings(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Client.UpdateSettings() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestClient_ReplaceSettings(t *testing.T) {
	server := newMockServer(mockServerExpectedInput{
		method:      http.MethodPut,
		path:        "/settings",
		contentType: "application/json",
		body:        `{"messages":{"limit_period":"PerHour","limit_value":50}}`,
	}, mockServerOutput{
		code: http.StatusOK,
		body: `{"messages":{"limit_period":"PerHour","limit_value":50}}`,
	})
	defer server.Close()

	client := newClient(server.URL)

	limitPeriod := smsgateway.PerHour
	limitValue := 50
	tests := []struct {
		name     string
		args     smsgateway.DeviceSettings
		expected smsgateway.DeviceSettings
		wantErr  bool
	}{
		{
			name: "Success",
			args: smsgateway.DeviceSettings{
				Messages: &smsgateway.SettingsMessages{
					LimitPeriod: &limitPeriod,
					LimitValue:  &limitValue,
				},
			},
			expected: smsgateway.DeviceSettings{
				Messages: &smsgateway.SettingsMessages{
					LimitPeriod: &limitPeriod,
					LimitValue:  &limitValue,
				},
			},
			wantErr: false,
		},
		{
			name:     "Error",
			args:     smsgateway.DeviceSettings{},
			expected: smsgateway.DeviceSettings{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.ReplaceSettings(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Client.ReplaceSettings() got = %v, want %v", got, tt.expected)
			}
		})
	}
}
