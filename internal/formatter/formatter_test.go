package formatter

import (
	"testing"

	"github.com/putcho01/atlassian-cli/internal/jira"
)

func TestParseFormat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  Format
	}{
		{"lowercase json", "json", JSON},
		{"uppercase JSON", "JSON", JSON},
		{"mixed case Json", "Json", JSON},
		{"lowercase markdown", "markdown", Markdown},
		{"alias md", "md", Markdown},
		{"uppercase MD", "MD", Markdown},
		{"table explicit", "table", Table},
		{"empty defaults to table", "", Table},
		{"unknown defaults to table", "csv", Table},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ParseFormat(tt.input)
			if got != tt.want {
				t.Errorf("ParseFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSafeGetters(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"safeName nil", safeName(nil), ""},
		{"safeName non-nil", safeName(&jira.Status{Name: "Done"}), "Done"},
		{"safeIssueType nil", safeIssueType(nil), ""},
		{"safeIssueType non-nil", safeIssueType(&jira.IssueType{Name: "Story"}), "Story"},
		{"safePriority nil", safePriority(nil), ""},
		{"safePriority non-nil", safePriority(&jira.Priority{Name: "High"}), "High"},
		{"safeUser nil", safeUser(nil), "Unassigned"},
		{"safeUser non-nil", safeUser(&jira.User{DisplayName: "Alice"}), "Alice"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got != tt.want {
				t.Errorf("got %q, want %q", tt.got, tt.want)
			}
		})
	}
}
