package cmd

import (
	"github.com/putcho01/atlassian-cli/internal/jira"
	"github.com/spf13/cobra"
)

var jiraIssueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Jira issue",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newJiraClient()
		if err != nil {
			return err
		}

		project, _ := cmd.Flags().GetString("project")
		summary, _ := cmd.Flags().GetString("summary")
		issueType, _ := cmd.Flags().GetString("type")
		description, _ := cmd.Flags().GetString("description")
		assignee, _ := cmd.Flags().GetString("assignee")
		priority, _ := cmd.Flags().GetString("priority")
		labels, _ := cmd.Flags().GetStringSlice("labels")

		issue, err := client.CreateIssue(cmd.Context(), &jira.CreateIssueInput{
			ProjectKey:  project,
			Summary:     summary,
			IssueType:   issueType,
			Description: description,
			Assignee:    assignee,
			Priority:    priority,
			Labels:      labels,
		})
		if err != nil {
			return err
		}
		return newFormatter(cmd).PrintIssue(issue)
	},
}

func init() {
	jiraIssueCreateCmd.Flags().String("project", "", "Project key (required)")
	jiraIssueCreateCmd.Flags().String("summary", "", "Issue summary (required)")
	jiraIssueCreateCmd.Flags().String("type", "Task", "Issue type")
	jiraIssueCreateCmd.Flags().String("description", "", "Issue description")
	jiraIssueCreateCmd.Flags().String("assignee", "", "Assignee account ID")
	jiraIssueCreateCmd.Flags().String("priority", "", "Priority name")
	jiraIssueCreateCmd.Flags().StringSlice("labels", nil, "Labels")
	_ = jiraIssueCreateCmd.MarkFlagRequired("project")
	_ = jiraIssueCreateCmd.MarkFlagRequired("summary")
	jiraIssueCmd.AddCommand(jiraIssueCreateCmd)
}
