package formatter

import (
	"bytes"
	"strings"
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
		{"safeUserDisplay nil unknown", safeUserDisplay(nil, "Unknown"), "Unknown"},
		{"safeUserDisplay nil unassigned", safeUserDisplay(nil, "Unassigned"), "Unassigned"},
		{"safeUserDisplay non-nil", safeUserDisplay(&jira.User{DisplayName: "Bob"}, "Unknown"), "Bob"},
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

func TestPrintComments(t *testing.T) {
	t.Parallel()
	comments := []jira.Comment{
		{ID: "1", Author: &jira.User{DisplayName: "Alice"}, Body: "Hello", Created: "2024-01-01"},
		{ID: "2", Author: nil, Body: "World", Created: "2024-01-02"},
	}

	tests := []struct {
		name     string
		format   Format
		contains []string
	}{
		{
			name:     "table contains headers and IDs",
			format:   Table,
			contains: []string{"ID", "AUTHOR", "CREATED", "1", "Alice", "Hello", "2", "Unknown"},
		},
		{
			name:     "markdown contains author and body",
			format:   Markdown,
			contains: []string{"### Alice", "Hello", "### Unknown", "World"},
		},
		{
			name:     "json contains comment fields",
			format:   JSON,
			contains: []string{`"id"`, `"1"`, `"Alice"`},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			f := New(tt.format, &buf)
			if err := f.PrintComments(comments); err != nil {
				t.Fatalf("PrintComments() error: %v", err)
			}
			out := buf.String()
			for _, want := range tt.contains {
				if !strings.Contains(out, want) {
					t.Errorf("output does not contain %q\ngot: %s", want, out)
				}
			}
		})
	}
}

func TestPrintComment(t *testing.T) {
	t.Parallel()
	comment := &jira.Comment{
		ID:      "10",
		Author:  &jira.User{DisplayName: "Alice"},
		Body:    "Single comment",
		Created: "2024-01-01",
	}

	tests := []struct {
		name     string
		format   Format
		contains []string
	}{
		{
			name:     "table contains ID and author",
			format:   Table,
			contains: []string{"10", "Alice", "2024-01-01", "Single comment"},
		},
		{
			name:     "markdown contains heading and body",
			format:   Markdown,
			contains: []string{"### Alice", "Single comment"},
		},
		{
			name:     "json contains id field",
			format:   JSON,
			contains: []string{`"id"`, `"10"`},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			f := New(tt.format, &buf)
			if err := f.PrintComment(comment); err != nil {
				t.Fatalf("PrintComment() error: %v", err)
			}
			out := buf.String()
			for _, want := range tt.contains {
				if !strings.Contains(out, want) {
					t.Errorf("output does not contain %q\ngot: %s", want, out)
				}
			}
		})
	}
}
