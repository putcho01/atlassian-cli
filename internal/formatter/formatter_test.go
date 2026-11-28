package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/putcho01/atlassian-cli/internal/confluence"
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

func TestMarkdownKeepsJiraTextVerbatim(t *testing.T) {
	t.Parallel()
	const text = "use List<String> in {code}Map<K,V>{code}"

	var buf bytes.Buffer
	f := New(Markdown, &buf)
	if err := f.PrintIssue(&jira.Issue{Key: "P-1", Fields: jira.IssueFields{Description: text}}); err != nil {
		t.Fatalf("PrintIssue() error: %v", err)
	}
	if err := f.PrintComment(&jira.Comment{Body: text}); err != nil {
		t.Fatalf("PrintComment() error: %v", err)
	}
	if got := strings.Count(buf.String(), text); got != 2 {
		t.Errorf("verbatim text appears %d times, want 2\ngot: %s", got, buf.String())
	}

	buf.Reset()
	if err := New(Table, &buf).PrintComment(&jira.Comment{Body: text}); err != nil {
		t.Fatalf("PrintComment() error: %v", err)
	}
	if !strings.Contains(buf.String(), text) {
		t.Errorf("table output does not contain %q\ngot: %s", text, buf.String())
	}
}

func TestMarkdownTablesEscapeCells(t *testing.T) {
	t.Parallel()
	user := &jira.User{DisplayName: "a|b"}
	tests := []struct {
		name  string
		print func(*Formatter) error
	}{
		{"user", func(f *Formatter) error { return f.PrintUser(user) }},
		{"issue", func(f *Formatter) error {
			return f.PrintIssue(&jira.Issue{Fields: jira.IssueFields{Labels: []string{"a|b"}}})
		}},
		{"issue list", func(f *Formatter) error {
			return f.PrintIssueList([]jira.Issue{{Fields: jira.IssueFields{Summary: "a|b"}}})
		}},
		{"transitions", func(f *Formatter) error { return f.PrintTransitions([]jira.Transition{{Name: "a|b"}}) }},
		{"page", func(f *Formatter) error { return f.PrintPage(&confluence.Page{Status: "a|b"}) }},
		{"labels", func(f *Formatter) error { return f.PrintLabels([]confluence.Label{{Name: "a|b"}}) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var buf bytes.Buffer
			if err := tt.print(New(Markdown, &buf)); err != nil {
				t.Fatalf("print error: %v", err)
			}
			if !strings.Contains(buf.String(), `| a\|b |`) {
				t.Errorf("cell is not escaped\ngot: %s", buf.String())
			}
		})
	}
}
