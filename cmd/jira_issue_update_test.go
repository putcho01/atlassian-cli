package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestJiraIssueUpdateFieldKeepsCommas(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode request: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	t.Setenv("JIRA_URL", srv.URL)
	t.Setenv("JIRA_EMAIL", "")
	t.Setenv("JIRA_API_TOKEN", "test-token")

	rootCmd.SetArgs([]string{
		"jira", "issue", "update", "P-1",
		"--field", "summary=Hello, world",
		"--field", `labels=["a","b"]`,
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	want := map[string]any{"fields": map[string]any{
		"summary": "Hello, world",
		"labels":  []any{"a", "b"},
	}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("request body diff (-want +got):\n%s", diff)
	}
}
