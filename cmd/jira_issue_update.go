package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var jiraIssueUpdateCmd = &cobra.Command{
	Use:   "update <issue-key>",
	Short: "Update a Jira issue",
	Long:  "Update a Jira issue field. Use --field key=value to set fields.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newJiraClient()
		if err != nil {
			return err
		}

		fieldFlags, _ := cmd.Flags().GetStringSlice("field")
		fields := make(map[string]any)
		for _, f := range fieldFlags {
			parts := strings.SplitN(f, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid field format %q, expected key=value", f)
			}
			var val any
			if err := json.Unmarshal([]byte(parts[1]), &val); err != nil {
				val = parts[1]
			}
			fields[parts[0]] = val
		}

		if len(fields) == 0 {
			return fmt.Errorf("at least one --field is required")
		}

		if err := client.UpdateIssue(cmd.Context(), args[0], fields); err != nil {
			return err
		}
		return newFormatter(cmd).PrintMessage(fmt.Sprintf("Issue %s updated successfully", args[0]))
	},
}

func init() {
	jiraIssueUpdateCmd.Flags().StringSlice("field", nil, "Field to update (key=value)")
	jiraIssueCmd.AddCommand(jiraIssueUpdateCmd)
}
