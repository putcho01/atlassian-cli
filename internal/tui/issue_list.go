package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/putcho01/atlassian-cli/internal/jira"
)

var (
	styleSelected  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	styleKey       = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Width(12)
	styleStatus    = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Width(16)
	styleSummary   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	styleHelp      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	styleDetail    = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	styleSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type Action int

const (
	ActionNone Action = iota
	ActionOpen
)

type Result struct {
	Issue  *jira.Issue
	Action Action
}

type issueListModel struct {
	issues    []jira.Issue
	cursor    int
	detail    bool
	result    *Result
	quitting  bool
	separator string
}

var keys = struct {
	up, down, enter, open, back, quit key.Binding
}{
	up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "detail")),
	open:  key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open in browser")),
	back:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}

func (m issueListModel) Init() tea.Cmd { return nil }

func (m issueListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		width := msg.Width
		if width > 80 {
			width = 80
		}
		m.separator = styleSeparator.Render(strings.Repeat("─", width))

	case tea.KeyMsg:
		if m.detail {
			switch {
			case key.Matches(msg, keys.open):
				return m.selectOpen()
			case key.Matches(msg, keys.back):
				m.detail = false
			case key.Matches(msg, keys.quit):
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
		}

		switch {
		case key.Matches(msg, keys.up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, keys.down):
			if m.cursor < len(m.issues)-1 {
				m.cursor++
			}
		case key.Matches(msg, keys.enter):
			m.detail = true
		case key.Matches(msg, keys.open):
			return m.selectOpen()
		case key.Matches(msg, keys.quit):
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m issueListModel) selectOpen() (tea.Model, tea.Cmd) {
	m.result = &Result{Issue: &m.issues[m.cursor], Action: ActionOpen}
	m.quitting = true
	return m, tea.Quit
}

func (m issueListModel) View() string {
	if m.quitting {
		return ""
	}
	if m.detail {
		return m.detailView()
	}
	return m.listView()
}

func (m issueListModel) listView() string {
	var sb strings.Builder
	sb.WriteString(styleSelected.Render(fmt.Sprintf("  Issues (%d)", len(m.issues))) + "\n")
	sb.WriteString(m.separator + "\n")

	for i, issue := range m.issues {
		issueKey := styleKey.Render(issue.Key)
		status := styleStatus.Render(safeName(issue.Fields.Status))
		summary := styleSummary.Render(truncate(issue.Fields.Summary, 50))

		row := fmt.Sprintf("  %s %s %s", issueKey, status, summary)
		if i == m.cursor {
			row = styleSelected.Render("▶ ") + issueKey + " " + status + " " + summary
		}
		sb.WriteString(row + "\n")
	}

	sb.WriteString("\n")
	sb.WriteString(styleHelp.Render("  ↑/↓ move  enter detail  o open in browser  q quit"))
	return sb.String()
}

func (m issueListModel) detailView() string {
	issue := &m.issues[m.cursor]
	var sb strings.Builder

	sb.WriteString(styleSelected.Render(fmt.Sprintf("  %s: %s", issue.Key, issue.Fields.Summary)) + "\n")
	sb.WriteString(m.separator + "\n")

	rows := []struct{ label, value string }{
		{"Status", safeName(issue.Fields.Status)},
		{"Type", safeType(issue.Fields.IssueType)},
		{"Priority", safePriority(issue.Fields.Priority)},
		{"Assignee", safeUser(issue.Fields.Assignee)},
		{"Reporter", safeUser(issue.Fields.Reporter)},
		{"Created", issue.Fields.Created},
		{"Updated", issue.Fields.Updated},
	}
	for _, r := range rows {
		sb.WriteString(fmt.Sprintf("  %s %s\n", styleStatus.Render(r.label), styleDetail.Render(r.value)))
	}

	if issue.Fields.Description != "" {
		sb.WriteString("\n" + styleStatus.Render("Description") + "\n")
		sb.WriteString(styleDetail.Render(truncate(issue.Fields.Description, 400)) + "\n")
	}

	sb.WriteString("\n")
	sb.WriteString(styleHelp.Render("  o open in browser  esc back  q quit"))
	return sb.String()
}

// RunIssueList launches the TUI and returns the selected result (nil if quit without selecting an action).
func RunIssueList(issues []jira.Issue) (*Result, error) {
	m := issueListModel{issues: issues}
	p := tea.NewProgram(m, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return nil, err
	}
	return final.(issueListModel).result, nil
}

func safeName(s *jira.Status) string {
	if s == nil {
		return ""
	}
	return s.Name
}

func safeType(t *jira.IssueType) string {
	if t == nil {
		return ""
	}
	return t.Name
}

func safePriority(p *jira.Priority) string {
	if p == nil {
		return ""
	}
	return p.Name
}

func safeUser(u *jira.User) string {
	if u == nil {
		return "Unassigned"
	}
	return u.DisplayName
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-3]) + "..."
}
