package cmd

import (
	"github.com/spf13/cobra"
)

var jiraIssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Jira issue commands",
}

func init() {
	jiraCmd.AddCommand(jiraIssueCmd)
}
