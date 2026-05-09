package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/putcho01/atlassian-cli/internal/httpclient"
)

func newTestJiraClient(baseURL string) *Client {
	return &Client{http: httpclient.NewFromConfig(baseURL, "test@test.com", "token")}
}

func TestBuildCreateFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    *CreateIssueInput
		expected map[string]any
	}{
		{
			name: "minimal required fields only",
			input: &CreateIssueInput{
				ProjectKey: "PROJ",
				Summary:    "Fix bug",
				IssueType:  "Bug",
			},
			expected: map[string]any{
				"project":   map[string]string{"key": "PROJ"},
				"summary":   "Fix bug",
				"issuetype": map[string]string{"name": "Bug"},
			},
		},
		{
			name: "all optional fields populated",
			input: &CreateIssueInput{
				ProjectKey:  "PROJ",
				Summary:     "New Feature",
				IssueType:   "Story",
				Description: "As a user...",
				Assignee:    "acc123",
				Priority:    "High",
				Labels:      []string{"backend", "v2"},
				CustomFields: map[string]any{
					"customfield_10001": "value1",
				},
			},
			expected: map[string]any{
				"project":           map[string]string{"key": "PROJ"},
				"summary":           "New Feature",
				"issuetype":         map[string]string{"name": "Story"},
				"description":       "As a user...",
				"assignee":          map[string]string{"accountId": "acc123"},
				"priority":          map[string]string{"name": "High"},
				"labels":            []string{"backend", "v2"},
				"customfield_10001": "value1",
			},
		},
		{
			name: "empty description not added",
			input: &CreateIssueInput{
				ProjectKey:  "PROJ",
				Summary:     "Task",
				IssueType:   "Task",
				Description: "",
			},
			expected: map[string]any{
				"project":   map[string]string{"key": "PROJ"},
				"summary":   "Task",
				"issuetype": map[string]string{"name": "Task"},
			},
		},
		{
			name: "empty assignee not added",
			input: &CreateIssueInput{
				ProjectKey: "PROJ",
				Summary:    "Task",
				IssueType:  "Task",
				Assignee:   "",
			},
			expected: map[string]any{
				"project":   map[string]string{"key": "PROJ"},
				"summary":   "Task",
				"issuetype": map[string]string{"name": "Task"},
			},
		},
		{
			name: "empty labels not added",
			input: &CreateIssueInput{
				ProjectKey: "PROJ",
				Summary:    "Task",
				IssueType:  "Task",
				Labels:     []string{},
			},
			expected: map[string]any{
				"project":   map[string]string{"key": "PROJ"},
				"summary":   "Task",
				"issuetype": map[string]string{"name": "Task"},
			},
		},
		{
			name: "custom fields merged into result",
			input: &CreateIssueInput{
				ProjectKey: "PROJ",
				Summary:    "Task",
				IssueType:  "Task",
				CustomFields: map[string]any{
					"customfield_10000": 42,
				},
			},
			expected: map[string]any{
				"project":           map[string]string{"key": "PROJ"},
				"summary":           "Task",
				"issuetype":         map[string]string{"name": "Task"},
				"customfield_10000": 42,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := buildCreateFields(tt.input)
			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("buildCreateFields() diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTransitionIssue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		statusName     string
		transitions    []Transition
		wantErr        bool
		wantErrContain string
	}{
		{
			name:       "success by transition name",
			statusName: "In Progress",
			transitions: []Transition{
				{ID: "11", Name: "In Progress", To: &Status{Name: "In Progress"}},
				{ID: "21", Name: "Done", To: &Status{Name: "Done"}},
			},
			wantErr: false,
		},
		{
			name:       "success by To.Name match",
			statusName: "Done",
			transitions: []Transition{
				{ID: "21", Name: "Close Issue", To: &Status{Name: "Done"}},
			},
			wantErr: false,
		},
		{
			name:       "case insensitive match",
			statusName: "in progress",
			transitions: []Transition{
				{ID: "11", Name: "In Progress", To: &Status{Name: "In Progress"}},
			},
			wantErr: false,
		},
		{
			name:       "transition not found",
			statusName: "Nonexistent",
			transitions: []Transition{
				{ID: "11", Name: "In Progress", To: &Status{Name: "In Progress"}},
			},
			wantErr:        true,
			wantErrContain: "not found",
		},
		{
			name:           "empty transitions list",
			statusName:     "Done",
			transitions:    []Transition{},
			wantErr:        true,
			wantErrContain: "not found",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mux := http.NewServeMux()
			mux.HandleFunc("/rest/api/2/issue/TEST-1/transitions", func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(transitionsResponse{Transitions: tt.transitions})
				case http.MethodPost:
					w.WriteHeader(http.StatusNoContent)
				default:
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			})
			srv := httptest.NewServer(mux)
			defer srv.Close()

			err := newTestJiraClient(srv.URL).TransitionIssue(context.Background(), "TEST-1", tt.statusName)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if tt.wantErrContain != "" && !strings.Contains(err.Error(), tt.wantErrContain) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.wantErrContain)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
