package ca

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/android-sms-gateway/client-go/rest"
)

type Client struct {
	*rest.Client
}

// PostCSR posts a Certificate Signing Request (CSR) to the Certificate Authority (CA) service.
//
// The service will validate the CSR and respond with a request ID.
//
// The request ID can be used to get the status of the request using the GetCSRStatus method.
func (c *Client) PostCSR(ctx context.Context, request PostCSRRequest) (PostCSRResponse, error) {
	path := "/csr"
	resp := new(PostCSRResponse)

	if err := c.Do(ctx, http.MethodPost, path, emptyHeaders, &request, resp); err != nil {
		return *resp, fmt.Errorf("failed to post CSR: %w", err)
	}

	return *resp, nil
}

// GetCSRStatus retrieves the status of a Certificate Signing Request (CSR) from the Certificate Authority (CA) service.
func (c *Client) GetCSRStatus(ctx context.Context, requestID string) (GetCSRStatusResponse, error) {
	path := "/csr/" + url.PathEscape(requestID)
	resp := new(GetCSRStatusResponse)

	if err := c.Do(ctx, http.MethodGet, path, emptyHeaders, nil, resp); err != nil {
		return *resp, fmt.Errorf("failed to get CSR status: %w", err)
	}

	return *resp, nil
}

// NewClient creates a new instance of the CA API Client.
func NewClient(options ...Option) *Client {
	config := new(Config)
	for _, option := range options {
		option(config)
	}

	return &Client{
		Client: rest.NewClient(rest.Config{
			Client:  config.Client(),
			BaseURL: config.BaseURL(),
		}),
	}
}
