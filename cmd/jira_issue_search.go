package cmd

import (
	"github.com/putcho01/atlassian-cli/internal/tui"
	"github.com/spf13/cobra"
)

var jiraIssueSearchCmd = &cobra.Command{
	Use:   "search <jql>",
	Short: "Search Jira issues using JQL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newJiraClient()
		if err != nil {
			return err
		}
		maxResults, _ := cmd.Flags().GetInt("max-results")
		result, err := client.SearchIssues(cmd.Context(), args[0], maxResults)
		if err != nil {
			return err
		}

		interactive, _ := cmd.Flags().GetBool("interactive")
		if !interactive {
			return newFormatter(cmd).PrintIssueList(result.Issues)
		}

		res, err := tui.RunIssueList(result.Issues)
		if err != nil {
			return err
		}
		switch {
		case res == nil, res.Action == tui.ActionNone:
			return nil
		case res.Action == tui.ActionOpen:
			return openBrowser(issueURL(client.BaseURL(), res.Issue.Key))
		}
		return nil
	},
}

func init() {
	jiraIssueSearchCmd.Flags().Int("max-results", 50, "Maximum number of results")
	jiraIssueSearchCmd.Flags().BoolP("interactive", "i", false, "Launch interactive TUI picker")
	jiraIssueCmd.AddCommand(jiraIssueSearchCmd)
}
