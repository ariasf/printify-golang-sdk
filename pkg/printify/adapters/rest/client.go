// Package rest implements the driven ports against the Printify HTTP API.
// It owns all transport concerns; no business logic lives here.
package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/printify-go/pkg/printify/domain"
)

// DefaultBaseURL is the production Printify API endpoint.
const DefaultBaseURL = "https://api.printify.com"

// UserAgentVersion is the SDK version embedded in the User-Agent header of
// every request. It is injected at build time via:
//
//	go build -ldflags "-X github.com/printify-go/pkg/printify/adapters/rest.UserAgentVersion=vX.Y.Z"
var UserAgentVersion = "v0.0.1-alpha"

// Doer abstracts *http.Client for testability.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// AuthError indicates a 401/403 response from the API.
type AuthError struct {
	Code    int
	Message string
}

func (e *AuthError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("unauthorized: %s", e.Message)
	}
	return "unauthorized"
}

// Unwrap makes errors.Is(err, domain.ErrUnauthorized) work.
func (e *AuthError) Unwrap() error { return domain.ErrUnauthorized }

// Client is the shared HTTP transport for all REST gateways.
type Client struct {
	baseURL string
	token   string
	http    Doer
}

// NewClient builds a transport client. If doer is nil, http.DefaultClient is used.
func NewClient(baseURL, token string, doer Doer) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	if doer == nil {
		doer = http.DefaultClient
	}
	return &Client{baseURL: baseURL, token: token, http: doer}
}

// do performs a JSON request, optionally with query params and a request body,
// decoding the response into out (if non-nil).
func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encoding %s %s body: %w", method, path, err)
		}
		reader = bytes.NewReader(data)
	}

	target := c.baseURL + path
	if len(query) > 0 {
		target += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, target, reader)
	if err != nil {
		return fmt.Errorf("building request %s %s: %w", method, path, err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "printify-golang-sdk/"+UserAgentVersion)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return c.errorFromResponse(resp.StatusCode, resp.Body)
	}
	if out == nil {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil && err != io.EOF {
		return fmt.Errorf("decoding %s %s response: %w", method, path, err)
	}
	return nil
}

func (c *Client) errorFromResponse(status int, body io.Reader) error {
	var payload struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	_ = json.NewDecoder(body).Decode(&payload)

	switch status {
	case http.StatusUnauthorized, http.StatusForbidden:
		return &AuthError{Code: payload.Code, Message: payload.Message}
	case http.StatusNotFound:
		return domain.ErrNotFound
	case http.StatusConflict:
		return domain.ErrConflict
	default:
		if payload.Message != "" {
			return fmt.Errorf("printify api error %d: %s", status, payload.Message)
		}
		return fmt.Errorf("unexpected status %d", status)
	}
}
