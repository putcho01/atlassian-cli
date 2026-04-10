package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var jiraIssueTransitionCmd = &cobra.Command{
	Use:   "transition <issue-key> <status>",
	Short: "Transition a Jira issue to a new status",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newJiraClient()
		if err != nil {
			return err
		}
		if err := client.TransitionIssue(cmd.Context(), args[0], args[1]); err != nil {
			return err
		}
		return newFormatter(cmd).PrintMessage(fmt.Sprintf("Issue %s transitioned to %s", args[0], args[1]))
	},
}

func init() {
	jiraIssueCmd.AddCommand(jiraIssueTransitionCmd)
}
