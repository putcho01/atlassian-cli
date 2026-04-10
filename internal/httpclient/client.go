package httpclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	authHeader string
	httpClient *http.Client
}

type APIError struct {
	StatusCode int
	Messages   []string          `json:"errorMessages"`
	Errors     map[string]string `json:"errors"`
}

func (e *APIError) Error() string {
	parts := make([]string, 0, len(e.Messages)+len(e.Errors))
	parts = append(parts, e.Messages...)
	for k, v := range e.Errors {
		parts = append(parts, fmt.Sprintf("%s: %s", k, v))
	}
	if len(parts) == 0 {
		return fmt.Sprintf("API error: HTTP %d", e.StatusCode)
	}
	return fmt.Sprintf("API error (HTTP %d): %s", e.StatusCode, strings.Join(parts, "; "))
}

// NewFromConfig creates a client, choosing Basic auth (Cloud) when email is set,
// or Bearer auth (Server/DC) otherwise.
func NewFromConfig(baseURL, email, token string) *Client {
	var authHeader string
	if email != "" {
		cred := base64.StdEncoding.EncodeToString([]byte(email + ":" + token))
		authHeader = "Basic " + cred
	} else {
		authHeader = "Bearer " + token
	}
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		authHeader: authHeader,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Get performs a GET request and returns the raw response.
func (c *Client) Get(ctx context.Context, path string, query url.Values) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, path, query, nil)
}

// Do performs an HTTP request, checks the status, and decodes JSON into result.
func (c *Client) Do(ctx context.Context, method, path string, query url.Values, reqBody any, result any) error {
	resp, err := c.do(ctx, method, path, query, reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseAPIError(resp)
	}

	if result != nil && resp.StatusCode != http.StatusNoContent {
		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			preview, _ := io.ReadAll(io.LimitReader(resp.Body, 201))
			return fmt.Errorf("expected JSON response but got %q from %s %s\n%s", contentType, method, path, preview)
		}
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body any) (*http.Response, error) {
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

func parseAPIError(resp *http.Response) error {
	apiErr := &APIError{StatusCode: resp.StatusCode}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiErr
	}
	_ = json.Unmarshal(data, apiErr)
	return apiErr
}
