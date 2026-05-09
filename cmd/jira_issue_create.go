package cmd

import (
	"fmt"

	"github.com/putcho01/atlassian-cli/internal/config"
	"github.com/putcho01/atlassian-cli/internal/jira"
	"github.com/spf13/cobra"
)

var jiraIssueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Jira issue",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadJiraConfig()
		if err != nil {
			return err
		}
		client := jira.NewClient(cfg)

		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			project = cfg.DefaultProject
		}
		if project == "" {
			return fmt.Errorf("--project is required (or set JIRA_DEFAULT_PROJECT)")
		}

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
	jiraIssueCreateCmd.Flags().String("project", "", "Project key (falls back to JIRA_DEFAULT_PROJECT)")
	jiraIssueCreateCmd.Flags().String("summary", "", "Issue summary (required)")
	jiraIssueCreateCmd.Flags().String("type", "Task", "Issue type")
	jiraIssueCreateCmd.Flags().String("description", "", "Issue description")
	jiraIssueCreateCmd.Flags().String("assignee", "", "Assignee account ID")
	jiraIssueCreateCmd.Flags().String("priority", "", "Priority name")
	jiraIssueCreateCmd.Flags().StringSlice("labels", nil, "Labels")
	_ = jiraIssueCreateCmd.MarkFlagRequired("summary")
	jiraIssueCmd.AddCommand(jiraIssueCreateCmd)
}
