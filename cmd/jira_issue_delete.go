package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var jiraIssueDeleteCmd = &cobra.Command{
	Use:   "delete <issue-key>",
	Short: "Delete a Jira issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newJiraClient()
		if err != nil {
			return err
		}
		if err := client.DeleteIssue(cmd.Context(), args[0]); err != nil {
			return err
		}
		return newFormatter(cmd).PrintMessage(fmt.Sprintf("Issue %s deleted", args[0]))
	},
}

func init() {
	jiraIssueCmd.AddCommand(jiraIssueDeleteCmd)
}
