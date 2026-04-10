package cmd

import (
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
		return newFormatter(cmd).PrintIssueList(result.Issues)
	},
}

func init() {
	jiraIssueSearchCmd.Flags().Int("max-results", 50, "Maximum number of results")
	jiraIssueCmd.AddCommand(jiraIssueSearchCmd)
}
