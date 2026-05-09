package cmd

import (
	"github.com/spf13/cobra"
)

var jiraIssueCommentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage comments on a Jira issue",
}

var jiraIssueCommentListCmd = &cobra.Command{
	Use:   "list <issue-key>",
	Short: "List comments on a Jira issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newJiraClient()
		if err != nil {
			return err
		}
		comments, err := client.ListComments(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return newFormatter(cmd).PrintComments(comments)
	},
}

var jiraIssueCommentAddCmd = &cobra.Command{
	Use:   "add <issue-key>",
	Short: "Add a comment to a Jira issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body, _ := cmd.Flags().GetString("body")
		client, err := newJiraClient()
		if err != nil {
			return err
		}
		comment, err := client.AddComment(cmd.Context(), args[0], body)
		if err != nil {
			return err
		}
		return newFormatter(cmd).PrintComment(comment)
	},
}

func init() {
	jiraIssueCommentAddCmd.Flags().String("body", "", "Comment body (required)")
	_ = jiraIssueCommentAddCmd.MarkFlagRequired("body")
	jiraIssueCommentCmd.AddCommand(jiraIssueCommentListCmd)
	jiraIssueCommentCmd.AddCommand(jiraIssueCommentAddCmd)
	jiraIssueCmd.AddCommand(jiraIssueCommentCmd)
}
