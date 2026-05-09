package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/putcho01/atlassian-cli/internal/confluence"
	"github.com/putcho01/atlassian-cli/internal/htmlconv"
	"github.com/putcho01/atlassian-cli/internal/jira"
)

type Format string

const (
	Table    Format = "table"
	JSON     Format = "json"
	Markdown Format = "markdown"
)

type Formatter struct {
	format Format
	writer io.Writer
}

func New(format Format, w io.Writer) *Formatter {
	if w == nil {
		w = os.Stdout
	}
	return &Formatter{format: format, writer: w}
}

func ParseFormat(s string) Format {
	switch strings.ToLower(s) {
	case "json":
		return JSON
	case "markdown", "md":
		return Markdown
	default:
		return Table
	}
}

// PrintJSON outputs any value as indented JSON.
func (f *Formatter) PrintJSON(v any) error {
	enc := json.NewEncoder(f.writer)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// PrintUser prints a Jira user.
func (f *Formatter) PrintUser(user *jira.User) error {
	switch f.format {
	case JSON:
		return f.PrintJSON(user)
	case Markdown:
		fmt.Fprintf(f.writer, "## User: %s\n\n", user.DisplayName)
		fmt.Fprintf(f.writer, "| Field | Value |\n|-------|-------|\n")
		fmt.Fprintf(f.writer, "| Account ID | %s |\n", user.AccountID)
		fmt.Fprintf(f.writer, "| Display Name | %s |\n", user.DisplayName)
		fmt.Fprintf(f.writer, "| Email | %s |\n", user.EmailAddress)
		fmt.Fprintf(f.writer, "| Active | %v |\n", user.Active)
		return nil
	default:
		tw := tabwriter.NewWriter(f.writer, 0, 4, 2, ' ', 0)
		fmt.Fprintf(tw, "ACCOUNT ID\tDISPLAY NAME\tEMAIL\tACTIVE\n")
		fmt.Fprintf(tw, "%s\t%s\t%s\t%v\n", user.AccountID, user.DisplayName, user.EmailAddress, user.Active)
		return tw.Flush()
	}
}

// PrintIssue prints a single Jira issue.
func (f *Formatter) PrintIssue(issue *jira.Issue) error {
	switch f.format {
	case JSON:
		return f.PrintJSON(issue)
	case Markdown:
		fmt.Fprintf(f.writer, "## %s: %s\n\n", issue.Key, issue.Fields.Summary)
		fmt.Fprintf(f.writer, "| Field | Value |\n|-------|-------|\n")
		fmt.Fprintf(f.writer, "| Status | %s |\n", safeName(issue.Fields.Status))
		fmt.Fprintf(f.writer, "| Type | %s |\n", safeIssueType(issue.Fields.IssueType))
		fmt.Fprintf(f.writer, "| Priority | %s |\n", safePriority(issue.Fields.Priority))
		fmt.Fprintf(f.writer, "| Assignee | %s |\n", safeUser(issue.Fields.Assignee))
		fmt.Fprintf(f.writer, "| Reporter | %s |\n", safeUser(issue.Fields.Reporter))
		fmt.Fprintf(f.writer, "| Labels | %s |\n", strings.Join(issue.Fields.Labels, ", "))
		fmt.Fprintf(f.writer, "| Created | %s |\n", issue.Fields.Created)
		fmt.Fprintf(f.writer, "| Updated | %s |\n", issue.Fields.Updated)
		if issue.Fields.Description != "" {
			fmt.Fprintf(f.writer, "\n### Description\n\n%s\n", htmlconv.Convert(issue.Fields.Description))
		}
		return nil
	default:
		tw := tabwriter.NewWriter(f.writer, 0, 4, 2, ' ', 0)
		fmt.Fprintf(tw, "KEY\tSUMMARY\tSTATUS\tASSIGNEE\tPRIORITY\n")
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
			issue.Key, issue.Fields.Summary,
			safeName(issue.Fields.Status),
			safeUser(issue.Fields.Assignee),
			safePriority(issue.Fields.Priority))
		return tw.Flush()
	}
}

// PrintIssueList prints a list of Jira issues.
func (f *Formatter) PrintIssueList(issues []jira.Issue) error {
	switch f.format {
	case JSON:
		return f.PrintJSON(issues)
	case Markdown:
		fmt.Fprintf(f.writer, "| Key | Summary | Status | Assignee | Priority |\n")
		fmt.Fprintf(f.writer, "|-----|---------|--------|----------|----------|\n")
		for _, issue := range issues {
			fmt.Fprintf(f.writer, "| %s | %s | %s | %s | %s |\n",
				issue.Key, issue.Fields.Summary,
				safeName(issue.Fields.Status),
				safeUser(issue.Fields.Assignee),
				safePriority(issue.Fields.Priority))
		}
		return nil
	default:
		tw := tabwriter.NewWriter(f.writer, 0, 4, 2, ' ', 0)
		fmt.Fprintf(tw, "KEY\tSUMMARY\tSTATUS\tASSIGNEE\tPRIORITY\n")
		for _, issue := range issues {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
				issue.Key, issue.Fields.Summary,
				safeName(issue.Fields.Status),
				safeUser(issue.Fields.Assignee),
				safePriority(issue.Fields.Priority))
		}
		return tw.Flush()
	}
}

// PrintTransitions prints available transitions.
func (f *Formatter) PrintTransitions(transitions []jira.Transition) error {
	switch f.format {
	case JSON:
		return f.PrintJSON(transitions)
	case Markdown:
		fmt.Fprintf(f.writer, "| ID | Name | To |\n|-----|------|----|\n")
		for _, t := range transitions {
			to := ""
			if t.To != nil {
				to = t.To.Name
			}
			fmt.Fprintf(f.writer, "| %s | %s | %s |\n", t.ID, t.Name, to)
		}
		return nil
	default:
		tw := tabwriter.NewWriter(f.writer, 0, 4, 2, ' ', 0)
		fmt.Fprintf(tw, "ID\tNAME\tTO\n")
		for _, t := range transitions {
			to := ""
			if t.To != nil {
				to = t.To.Name
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\n", t.ID, t.Name, to)
		}
		return tw.Flush()
	}
}

// PrintPage prints a Confluence page.
func (f *Formatter) PrintPage(page *confluence.Page) error {
	switch f.format {
	case JSON:
		return f.PrintJSON(page)
	case Markdown:
		fmt.Fprintf(f.writer, "## %s\n\n", page.Title)
		fmt.Fprintf(f.writer, "| Field | Value |\n|-------|-------|\n")
		fmt.Fprintf(f.writer, "| ID | %s |\n", page.ID)
		fmt.Fprintf(f.writer, "| Status | %s |\n", page.Status)
		fmt.Fprintf(f.writer, "| Version | %d |\n", page.Version.Number)
		if page.Body.Storage != nil {
			fmt.Fprintf(f.writer, "\n### Content\n\n%s\n", htmlconv.Convert(page.Body.Storage.Value))
		}
		return nil
	default:
		tw := tabwriter.NewWriter(f.writer, 0, 4, 2, ' ', 0)
		fmt.Fprintf(tw, "ID\tTITLE\tSTATUS\tVERSION\n")
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\n", page.ID, page.Title, page.Status, page.Version.Number)
		tw.Flush()
		if page.Body.Storage != nil {
			fmt.Fprintf(f.writer, "\n%s\n", htmlconv.Convert(page.Body.Storage.Value))
		}
		return nil
	}
}

// PrintLabels prints Confluence labels.
func (f *Formatter) PrintLabels(labels []confluence.Label) error {
	switch f.format {
	case JSON:
		return f.PrintJSON(labels)
	case Markdown:
		fmt.Fprintf(f.writer, "| Name | Prefix |\n|------|--------|\n")
		for _, l := range labels {
			fmt.Fprintf(f.writer, "| %s | %s |\n", l.Name, l.Prefix)
		}
		return nil
	default:
		tw := tabwriter.NewWriter(f.writer, 0, 4, 2, ' ', 0)
		fmt.Fprintf(tw, "NAME\tPREFIX\n")
		for _, l := range labels {
			fmt.Fprintf(tw, "%s\t%s\n", l.Name, l.Prefix)
		}
		return tw.Flush()
	}
}

func (f *Formatter) PrintComments(comments []jira.Comment) error {
	switch f.format {
	case JSON:
		return f.PrintJSON(comments)
	case Markdown:
		for i := range comments {
			if err := f.PrintComment(&comments[i]); err != nil {
				return err
			}
		}
		return nil
	default:
		tw := tabwriter.NewWriter(f.writer, 0, 4, 2, ' ', 0)
		fmt.Fprintf(tw, "ID\tAUTHOR\tCREATED\tBODY\n")
		for _, c := range comments {
			body := []rune(c.Body)
			if len(body) > 80 {
				body = append(body[:77], []rune("...")...)
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", c.ID, safeUserDisplay(c.Author, "Unknown"), c.Created, string(body))
		}
		return tw.Flush()
	}
}

func (f *Formatter) PrintComment(c *jira.Comment) error {
	switch f.format {
	case JSON:
		return f.PrintJSON(c)
	case Markdown:
		fmt.Fprintf(f.writer, "### %s — %s\n\n%s\n", safeUserDisplay(c.Author, "Unknown"), c.Created, htmlconv.Convert(c.Body))
		return nil
	default:
		fmt.Fprintf(f.writer, "ID\tAUTHOR\tCREATED\n")
		fmt.Fprintf(f.writer, "%s\t%s\t%s\n", c.ID, safeUserDisplay(c.Author, "Unknown"), c.Created)
		fmt.Fprintf(f.writer, "\n%s\n", htmlconv.Convert(c.Body))
		return nil
	}
}

// PrintMessage prints a simple success/info message.
func (f *Formatter) PrintMessage(msg string) error {
	switch f.format {
	case JSON:
		return f.PrintJSON(map[string]string{"message": msg})
	default:
		_, err := fmt.Fprintln(f.writer, msg)
		return err
	}
}

func safeName(s *jira.Status) string {
	if s == nil {
		return ""
	}
	return s.Name
}

func safeIssueType(t *jira.IssueType) string {
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

func safeUserDisplay(u *jira.User, fallback string) string {
	if u == nil {
		return fallback
	}
	return u.DisplayName
}

func safeUser(u *jira.User) string {
	return safeUserDisplay(u, "Unassigned")
}
