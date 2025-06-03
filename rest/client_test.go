package rest_test

import (
	"context"
	"errors"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/android-sms-gateway/client-go/rest"
)

func setupTestServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			if r.Method != http.MethodGet {
				t.Errorf("Expected method GET, got %s", r.Method)
			}
		case "/204":
			w.WriteHeader(http.StatusNoContent)
			return
		case "/400":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("bad request"))
			return
		case "/404":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("not found"))
			return
		case "/409":
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte("conflict"))
			return
		case "/500":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("internal server error"))
			return
		case "/corrupt":
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "application/json")
			_, _ = w.Write([]byte("{not a json"))
			return
		}

		if r.URL.Path == "/body" {
			if r.Method != http.MethodPost {
				t.Errorf("Expected method POST, got %s", r.Method)
			}
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
			}

			_, _ = io.ReadAll(r.Body)
			defer r.Body.Close()
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id": "123", "state": "Pending"}`))
	}))
}

func TestClient_Do(t *testing.T) {
	httpServer := setupTestServer(t)
	defer httpServer.Close()

	type fields struct {
		config rest.Config
	}
	type args struct {
		ctx      context.Context
		method   string
		path     string
		headers  map[string]string
		payload  any
		response any
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantErrType error
	}{
		{
			name: "Empty method",
			fields: fields{
				config: rest.Config{},
			},
			args: args{
				ctx:    context.Background(),
				method: "",
				path:   "/",
			},
			wantErr:     true,
			wantErrType: nil, // No specific error type expected
		},
		{
			name: "With body",
			fields: fields{
				config: rest.Config{
					BaseURL: httpServer.URL,
				},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodPost,
				path:   "/body",
				payload: map[string]string{
					"foo": "bar",
				},
			},
			wantErr:     false,
			wantErrType: nil, // No error expected
		},
		{
			name: "HTTP 400 error",
			fields: fields{
				config: rest.Config{
					BaseURL: httpServer.URL,
				},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				path:   "/400",
			},
			wantErr:     true,
			wantErrType: rest.ErrBadRequest,
		},
		{
			name: "HTTP 404 error",
			fields: fields{
				config: rest.Config{
					BaseURL: httpServer.URL,
				},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				path:   "/404",
			},
			wantErr:     true,
			wantErrType: rest.ErrClient,
		},
		{
			name: "HTTP 409 error",
			fields: fields{
				config: rest.Config{
					BaseURL: httpServer.URL,
				},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				path:   "/409",
			},
			wantErr:     true,
			wantErrType: rest.ErrConflict,
		},
		{
			name: "HTTP 500 error",
			fields: fields{
				config: rest.Config{
					BaseURL: httpServer.URL,
				},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				path:   "/500",
			},
			wantErr:     true,
			wantErrType: rest.ErrServer,
		},
		{
			name: "No Content response",
			fields: fields{
				config: rest.Config{
					BaseURL: httpServer.URL,
				},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				path:   "/204",
			},
			wantErr:     false,
			wantErrType: nil, // No error expected
		},
		{
			name: "Corrupt response",
			fields: fields{
				config: rest.Config{
					BaseURL: httpServer.URL,
				},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				path:   "/corrupt",
			},
			wantErr:     true,
			wantErrType: nil, // No specific error type expected
		},
		{
			name: "Corrupt request",
			fields: fields{
				config: rest.Config{
					BaseURL: httpServer.URL,
				},
			},
			args: args{
				ctx:     context.Background(),
				method:  http.MethodPost,
				path:    "/",
				payload: math.NaN(),
			},
			wantErr:     true,
			wantErrType: nil, // No specific error type expected
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := rest.NewClient(tt.fields.config)
			err := c.Do(tt.args.ctx, tt.args.method, tt.args.path, tt.args.headers, tt.args.payload, tt.args.response)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Skip error type checking if we don't expect an error
			if !tt.wantErr {
				return
			}

			// Verify error type if expected
			if tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
				t.Errorf("Expected error of type %v, got %v", tt.wantErrType, err)
			}
		})
	}
}
