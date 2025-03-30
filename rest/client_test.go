package rest_test

import (
	"context"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/android-sms-gateway/client-go/rest"
)

func TestClient_Do(t *testing.T) {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/204" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.URL.Path == "/404" {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("not found"))
			return
		}

		if r.URL.Path == "/corrupt" {
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
		name    string
		fields  fields
		args    args
		wantErr bool
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
			wantErr: true,
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
			wantErr: false,
		},
		{
			name: "HTTP error",
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
			wantErr: true,
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
			wantErr: false,
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
			wantErr: true,
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := rest.NewClient(tt.fields.config)
			if err := c.Do(tt.args.ctx, tt.args.method, tt.args.path, tt.args.headers, tt.args.payload, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("Client.Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
