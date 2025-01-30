package ca_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/android-sms-gateway/client-go/ca"
)

func TestClient_PostCSR(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/csr" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		req, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		if string(req) != `{"content":"-----BEGIN CERTIFICATE REQUEST-----"}` {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"request_id":"123", "status":"pending", "message": "CSR submitted successfully. Await processing."}`))
	}))
	defer server.Close()

	client := ca.NewClient(ca.WithBaseURL(server.URL))

	tests := []struct {
		name    string
		req     ca.PostCSRRequest
		want    ca.PostCSRResponse
		wantErr bool
	}{
		{
			name: "Success",
			req: ca.PostCSRRequest{
				Content: "-----BEGIN CERTIFICATE REQUEST-----",
			},
			want: ca.PostCSRResponse{
				RequestID: "123",
				Status:    ca.CSRStatusPending,
				Message:   ca.CSRStatusDescriptionPending,
			},
			wantErr: false,
		},
		{
			name:    "Error",
			req:     ca.PostCSRRequest{},
			want:    ca.PostCSRResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.PostCSR(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PostCSR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.PostCSR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetCSRStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/csr/123" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		_, _ = io.ReadAll(r.Body)
		defer r.Body.Close()

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"request_id":"123", "status":"approved", "message": "CSR approved. The certificate is ready for download.", "certificate":"-----BEGIN CERTIFICATE-----"}`))
	}))
	defer server.Close()

	client := ca.NewClient(ca.WithBaseURL(server.URL))

	tests := []struct {
		name    string
		req     string
		want    ca.PostCSRResponse
		wantErr bool
	}{
		{
			name: "Success",
			req:  "123",
			want: ca.PostCSRResponse{
				RequestID:   "123",
				Status:      ca.CSRStatusApproved,
				Message:     ca.CSRStatusDescriptionApproved,
				Certificate: "-----BEGIN CERTIFICATE-----",
			},
			wantErr: false,
		},
		{
			name:    "Error",
			req:     "456",
			want:    ca.PostCSRResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetCSRStatus(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetCSRStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetCSRStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
