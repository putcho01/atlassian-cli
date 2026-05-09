package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/putcho01/atlassian-cli/internal/jira"
)

func testIssues() []jira.Issue {
	return []jira.Issue{
		{Key: "PROJ-1", Fields: jira.IssueFields{Summary: "First issue", Status: &jira.Status{Name: "Open"}}},
		{Key: "PROJ-2", Fields: jira.IssueFields{Summary: "Second issue", Status: &jira.Status{Name: "In Progress"}}},
		{Key: "PROJ-3", Fields: jira.IssueFields{Summary: "Third issue", Status: &jira.Status{Name: "Done"}}},
	}
}

func newModel() issueListModel {
	return issueListModel{issues: testIssues(), separator: strings.Repeat("─", 80)}
}

func sendKey(m issueListModel, k string) issueListModel {
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)})
	return next.(issueListModel)
}

func sendSpecialKey(m issueListModel, t tea.KeyType) issueListModel {
	next, _ := m.Update(tea.KeyMsg{Type: t})
	return next.(issueListModel)
}

func TestNavigation(t *testing.T) {
	t.Parallel()
	m := newModel()

	if m.cursor != 0 {
		t.Fatalf("initial cursor = %d, want 0", m.cursor)
	}

	m = sendSpecialKey(m, tea.KeyDown)
	if m.cursor != 1 {
		t.Errorf("after down cursor = %d, want 1", m.cursor)
	}

	m = sendKey(m, "j")
	if m.cursor != 2 {
		t.Errorf("after j cursor = %d, want 2", m.cursor)
	}

	// clamp at bottom
	m = sendSpecialKey(m, tea.KeyDown)
	if m.cursor != 2 {
		t.Errorf("cursor should clamp at bottom, got %d", m.cursor)
	}

	m = sendSpecialKey(m, tea.KeyUp)
	if m.cursor != 1 {
		t.Errorf("after up cursor = %d, want 1", m.cursor)
	}

	m = sendKey(m, "k")
	if m.cursor != 0 {
		t.Errorf("after k cursor = %d, want 0", m.cursor)
	}

	// clamp at top
	m = sendSpecialKey(m, tea.KeyUp)
	if m.cursor != 0 {
		t.Errorf("cursor should clamp at top, got %d", m.cursor)
	}
}

func TestDetailView(t *testing.T) {
	t.Parallel()
	m := newModel()

	// enter detail
	m = sendSpecialKey(m, tea.KeyEnter)
	if !m.detail {
		t.Fatal("expected detail=true after enter")
	}
	view := m.View()
	if !strings.Contains(view, "PROJ-1") {
		t.Errorf("detail view should contain issue key, got: %s", view)
	}

	// esc back to list
	m = sendSpecialKey(m, tea.KeyEsc)
	if m.detail {
		t.Fatal("expected detail=false after esc")
	}
}

func TestOpenAction(t *testing.T) {
	t.Parallel()
	m := newModel()
	m.cursor = 1

	m = sendKey(m, "o")
	if !m.quitting {
		t.Fatal("expected quitting=true after o")
	}
	if m.result == nil {
		t.Fatal("expected result to be set")
	}
	if m.result.Action != ActionOpen {
		t.Errorf("action = %v, want ActionOpen", m.result.Action)
	}
	if m.result.Issue.Key != "PROJ-2" {
		t.Errorf("issue key = %s, want PROJ-2", m.result.Issue.Key)
	}
}

func TestQuit(t *testing.T) {
	t.Parallel()
	m := newModel()
	m = sendKey(m, "q")
	if !m.quitting {
		t.Fatal("expected quitting=true after q")
	}
	if m.result != nil {
		t.Errorf("expected nil result on quit, got %+v", m.result)
	}
}

func TestListViewContents(t *testing.T) {
	t.Parallel()
	m := newModel()
	view := m.View()

	for _, issue := range testIssues() {
		if !strings.Contains(view, issue.Key) {
			t.Errorf("list view missing issue key %s", issue.Key)
		}
	}
	if !strings.Contains(view, "open in browser") {
		t.Errorf("list view missing help text")
	}
}

func TestTruncate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		n     int
		want  string
	}{
		{"hello", 10, "hello"},
		{"hello world", 8, "hello..."},
		{"日本語テスト", 5, "日本..."},
	}
	for _, tt := range tests {
		got := truncate(tt.input, tt.n)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.n, got, tt.want)
		}
	}
}
