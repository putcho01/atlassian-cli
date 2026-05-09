package httpclient

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestNewFromConfig(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		baseURL        string
		email          string
		token          string
		wantAuthHeader string
		wantBaseURL    string
	}{
		{
			name:           "cloud: email set uses Basic auth",
			baseURL:        "https://example.atlassian.net",
			email:          "user@example.com",
			token:          "mytoken",
			wantAuthHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte("user@example.com:mytoken")),
			wantBaseURL:    "https://example.atlassian.net",
		},
		{
			name:           "server/DC: no email uses Bearer auth",
			baseURL:        "https://jira.company.internal",
			email:          "",
			token:          "personaltoken",
			wantAuthHeader: "Bearer personaltoken",
			wantBaseURL:    "https://jira.company.internal",
		},
		{
			name:           "trailing slash in baseURL is removed",
			baseURL:        "https://example.com/",
			email:          "a@b.com",
			token:          "tok",
			wantAuthHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte("a@b.com:tok")),
			wantBaseURL:    "https://example.com",
		},
		{
			name:           "multiple trailing slashes all removed",
			baseURL:        "https://example.com///",
			email:          "",
			token:          "tok",
			wantAuthHeader: "Bearer tok",
			wantBaseURL:    "https://example.com",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := NewFromConfig(tt.baseURL, tt.email, tt.token)
			if c.authHeader != tt.wantAuthHeader {
				t.Errorf("authHeader = %q, want %q", c.authHeader, tt.wantAuthHeader)
			}
			if c.baseURL != tt.wantBaseURL {
				t.Errorf("baseURL = %q, want %q", c.baseURL, tt.wantBaseURL)
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		apiErr     *APIError
		wantPrefix string
		wantAll    []string
	}{
		{
			name: "only messages",
			apiErr: &APIError{
				StatusCode: 400,
				Messages:   []string{"Issue not found"},
			},
			wantPrefix: "API error (HTTP 400):",
			wantAll:    []string{"Issue not found"},
		},
		{
			name: "only errors map",
			apiErr: &APIError{
				StatusCode: 422,
				Errors:     map[string]string{"summary": "Field required"},
			},
			wantPrefix: "API error (HTTP 422):",
			wantAll:    []string{"summary: Field required"},
		},
		{
			name: "both empty returns generic message",
			apiErr: &APIError{
				StatusCode: 500,
			},
			wantPrefix: "API error: HTTP 500",
			wantAll:    nil,
		},
		{
			name: "messages and errors both present",
			apiErr: &APIError{
				StatusCode: 400,
				Messages:   []string{"Bad request"},
				Errors:     map[string]string{"field": "invalid"},
			},
			wantPrefix: "API error (HTTP 400):",
			wantAll:    []string{"Bad request", "field: invalid"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.apiErr.Error()
			if !strings.HasPrefix(got, tt.wantPrefix) {
				t.Errorf("Error() = %q, want prefix %q", got, tt.wantPrefix)
			}
			for _, s := range tt.wantAll {
				if !strings.Contains(got, s) {
					t.Errorf("Error() = %q, want it to contain %q", got, s)
				}
			}
		})
	}
}
