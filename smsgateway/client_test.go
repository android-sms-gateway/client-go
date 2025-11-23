package smsgateway_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	method        string
	path          string
	query         string
	authorization string
	contentType   string
	body          string
}

type mockServerOutput struct {
	code int
	body string
}

func newMockServer(input mockServerExpectedInput, output mockServerOutput) *httptest.Server {
	if input.authorization == "" {
		input.authorization = authorizationHeader
	}

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

		if r.Header.Get("Authorization") != input.authorization {
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

func newJWTClient(baseURL string) *smsgateway.Client {
	return smsgateway.NewClient(smsgateway.Config{
		BaseURL: baseURL,
		Token:   password,
	})
}

func TestJWTClient_Send(t *testing.T) {
	server := newMockServer(mockServerExpectedInput{
		method:        http.MethodPost,
		path:          "/messages",
		authorization: "Bearer " + password,
		contentType:   "application/json",
		body:          `{"textMessage":{"text":"Hello World!"},"phoneNumbers":["+1234567890"]}`,
	}, mockServerOutput{
		code: http.StatusCreated,
		body: `{}`,
	})
	defer server.Close()

	client := newJWTClient(server.URL)

	t.Run("Success", func(t *testing.T) {
		message := smsgateway.Message{
			TextMessage: &smsgateway.TextMessage{
				Text: "Hello World!",
			},
			PhoneNumbers: []string{"+1234567890"},
		}

		_, err := client.Send(context.Background(), message)
		if err != nil {
			t.Errorf("Send() error = %v", err)
		}
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
	tests := []struct {
		name    string
		code    int
		body    string
		want    smsgateway.HealthResponse
		wantErr bool
	}{
		{
			name:    "Success",
			code:    http.StatusOK,
			body:    `{"status": "ok"}`,
			want:    smsgateway.HealthResponse{Status: "ok"},
			wantErr: false,
		},
		{
			name:    "Error response",
			code:    http.StatusInternalServerError,
			body:    `{"error": "internal error"}`,
			want:    smsgateway.HealthResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newMockServer(mockServerExpectedInput{
				method: http.MethodGet,
				path:   "/health",
			}, mockServerOutput{
				code: tt.code,
				body: tt.body,
			})
			defer server.Close()

			client := newClient(server.URL)
			resp, err := client.CheckHealth(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckHealth error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(resp, tt.want) {
				t.Errorf("CheckHealth response = %v, want %v", resp, tt.want)
			}
		})
	}
}

func TestClient_ExportInbox(t *testing.T) {
	tests := []struct {
		name    string
		req     smsgateway.MessagesExportRequest
		code    int
		wantErr bool
	}{
		{
			name: "Success",
			req: smsgateway.MessagesExportRequest{
				DeviceID: "qTRWxZkF",
				Since:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				Until:    time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			code:    http.StatusOK,
			wantErr: false,
		},
		{
			name:    "Error response",
			req:     smsgateway.MessagesExportRequest{},
			code:    http.StatusInternalServerError,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newMockServer(mockServerExpectedInput{
				method:      http.MethodPost,
				path:        "/inbox/export",
				contentType: "application/json",
				body:        `{"deviceId":"qTRWxZkF","since":"2023-01-01T00:00:00Z","until":"2023-01-02T00:00:00Z"}`,
			}, mockServerOutput{
				code: tt.code,
				body: `{}`,
			})
			defer server.Close()

			client := newClient(server.URL)
			err := client.ExportInbox(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExportInbox error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetLogs(t *testing.T) {
	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		code    int
		body    string
		want    []smsgateway.LogEntry
		wantErr bool
	}{
		{
			name: "Success",
			code: http.StatusOK,
			body: `[{"createdAt":"2023-01-01T00:00:00Z","message":"test"}]`,
			want: []smsgateway.LogEntry{
				{
					Message:   "test",
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name:    "Error response",
			code:    http.StatusInternalServerError,
			body:    `{"error": "internal error"}`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newMockServer(mockServerExpectedInput{
				method: http.MethodGet,
				path:   "/logs",
				query: "from=" + url.QueryEscape(
					from.Format(time.RFC3339),
				) + "&to=" + url.QueryEscape(
					to.Format(time.RFC3339),
				),
			}, mockServerOutput{
				code: tt.code,
				body: tt.body,
			})
			defer server.Close()

			client := newClient(server.URL)
			logs, err := client.GetLogs(context.Background(), from, to)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLogs error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(logs, tt.want) {
				t.Errorf("GetLogs logs = %v, want %v", logs, tt.want)
			}
		})
	}
}

func TestClient_GetSettings(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		body    string
		want    smsgateway.DeviceSettings
		wantErr bool
	}{
		{
			name: "Success",
			code: http.StatusOK,
			body: `{"messages":{"log_lifetime_days":30}}`,
			want: smsgateway.DeviceSettings{
				Messages: &smsgateway.SettingsMessages{
					LogLifetimeDays: ptr(30),
				},
			},
			wantErr: false,
		},
		{
			name:    "Error response",
			code:    http.StatusInternalServerError,
			body:    `{"error": "internal error"}`,
			want:    smsgateway.DeviceSettings{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newMockServer(mockServerExpectedInput{
				method: http.MethodGet,
				path:   "/settings",
			}, mockServerOutput{
				code: tt.code,
				body: tt.body,
			})
			defer server.Close()

			client := newClient(server.URL)
			settings, err := client.GetSettings(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSettings error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(settings, tt.want) {
				t.Errorf("GetSettings settings = %v, want %v", settings, tt.want)
			}
		})
	}
}

func TestClient_UpdateSettings(t *testing.T) {
	tests := []struct {
		name     string
		settings smsgateway.DeviceSettings
		code     int
		body     string
		want     smsgateway.DeviceSettings
		wantErr  bool
	}{
		{
			name: "Success",
			settings: smsgateway.DeviceSettings{
				Messages: &smsgateway.SettingsMessages{
					LogLifetimeDays: ptr(30),
				},
			},
			code: http.StatusOK,
			body: `{"messages":{"log_lifetime_days":30}}`,
			want: smsgateway.DeviceSettings{
				Messages: &smsgateway.SettingsMessages{
					LogLifetimeDays: ptr(30),
				},
			},
			wantErr: false,
		},
		{
			name: "Error response",
			settings: smsgateway.DeviceSettings{
				Messages: &smsgateway.SettingsMessages{
					LogLifetimeDays: ptr(30),
				},
			},
			code:    http.StatusInternalServerError,
			body:    `{"error": "internal error"}`,
			want:    smsgateway.DeviceSettings{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newMockServer(mockServerExpectedInput{
				method:      http.MethodPatch,
				path:        "/settings",
				contentType: "application/json",
				body:        `{"messages":{"log_lifetime_days":30}}`,
			}, mockServerOutput{
				code: tt.code,
				body: tt.body,
			})
			defer server.Close()

			client := newClient(server.URL)
			resp, err := client.UpdateSettings(context.Background(), tt.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSettings error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(resp, tt.want) {
				t.Errorf("UpdateSettings response = %v, want %v", resp, tt.want)
			}
		})
	}
}

func TestClient_ReplaceSettings(t *testing.T) {
	tests := []struct {
		name     string
		settings smsgateway.DeviceSettings
		code     int
		body     string
		want     smsgateway.DeviceSettings
		wantErr  bool
	}{
		{
			name: "Success",
			settings: smsgateway.DeviceSettings{
				Messages: &smsgateway.SettingsMessages{
					LogLifetimeDays: ptr(30),
				},
			},
			code: http.StatusOK,
			body: `{"messages":{"log_lifetime_days":30}}`,
			want: smsgateway.DeviceSettings{
				Messages: &smsgateway.SettingsMessages{
					LogLifetimeDays: ptr(30),
				},
			},
			wantErr: false,
		},
		{
			name: "Error response",
			settings: smsgateway.DeviceSettings{
				Messages: &smsgateway.SettingsMessages{
					LogLifetimeDays: ptr(30),
				},
			},
			code:    http.StatusInternalServerError,
			body:    `{"error": "internal error"}`,
			want:    smsgateway.DeviceSettings{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newMockServer(mockServerExpectedInput{
				method:      http.MethodPut,
				path:        "/settings",
				contentType: "application/json",
				body:        `{"messages":{"log_lifetime_days":30}}`,
			}, mockServerOutput{
				code: tt.code,
				body: tt.body,
			})
			defer server.Close()

			client := newClient(server.URL)
			resp, err := client.ReplaceSettings(context.Background(), tt.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceSettings error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(resp, tt.want) {
				t.Errorf("ReplaceSettings response = %v, want %v", resp, tt.want)
			}
		})
	}
}

func TestClient_GenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		req     smsgateway.TokenRequest
		code    int
		body    string
		want    smsgateway.TokenResponse
		wantErr bool
	}{
		{
			name: "Success",
			req:  smsgateway.TokenRequest{Scopes: []string{"messages:read"}, TTL: 3600},
			code: http.StatusOK,
			body: `{"id":"token_id_example","token_type":"Bearer","access_token":"access_token_example","expires_at":"2025-01-01T00:00:00Z"}`,
			want: smsgateway.TokenResponse{
				ID:          "token_id_example",
				TokenType:   "Bearer",
				AccessToken: "access_token_example",
				ExpiresAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name:    "Error response",
			req:     smsgateway.TokenRequest{Scopes: []string{"messages:read"}, TTL: 3600},
			code:    http.StatusInternalServerError,
			body:    `{"error": "internal error"}`,
			want:    smsgateway.TokenResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newMockServer(mockServerExpectedInput{
				method:      http.MethodPost,
				path:        "/auth/token",
				contentType: "application/json",
				body:        `{"ttl":3600,"scopes":["messages:read"]}`,
			}, mockServerOutput{
				code: tt.code,
				body: tt.body,
			})
			defer server.Close()

			client := newClient(server.URL)
			resp, err := client.GenerateToken(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(resp, tt.want) {
				t.Errorf("GenerateToken response = %v, want %v", resp, tt.want)
			}
		})
	}
}

func TestClient_RevokeToken(t *testing.T) {
	tests := []struct {
		name    string
		jti     string
		code    int
		wantErr bool
	}{
		{
			name:    "Success",
			jti:     "abc123",
			code:    http.StatusNoContent,
			wantErr: false,
		},
		{
			name:    "Error response",
			jti:     "abc123",
			code:    http.StatusInternalServerError,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newMockServer(mockServerExpectedInput{
				method: http.MethodDelete,
				path:   "/auth/token/abc123",
			}, mockServerOutput{
				code: tt.code,
			})
			defer server.Close()

			client := newClient(server.URL)
			err := client.RevokeToken(context.Background(), tt.jti)
			if (err != nil) != tt.wantErr {
				t.Errorf("RevokeToken error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
