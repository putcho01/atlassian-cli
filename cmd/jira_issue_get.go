package cmd

import (
	"github.com/spf13/cobra"
)

var jiraIssueGetCmd = &cobra.Command{
	Use:   "get <issue-key>",
	Short: "Get a Jira issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newJiraClient()
		if err != nil {
			return err
		}
		issue, err := client.GetIssue(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return newFormatter(cmd).PrintIssue(issue)
	},
}

func init() {
	jiraIssueCmd.AddCommand(jiraIssueGetCmd)
}
