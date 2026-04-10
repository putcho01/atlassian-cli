package cmd

import (
	"github.com/spf13/cobra"
)

var jiraIssueSubtasksCmd = &cobra.Command{
	Use:   "subtasks <issue-key>",
	Short: "List subtasks of a Jira issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newJiraClient()
		if err != nil {
			return err
		}
		subtasks, err := client.GetSubtasks(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return newFormatter(cmd).PrintIssueList(subtasks)
	},
}

func init() {
	jiraIssueCmd.AddCommand(jiraIssueSubtasksCmd)
}
