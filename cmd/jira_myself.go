package cmd

import (
	"github.com/spf13/cobra"
)

var jiraMyselfCmd = &cobra.Command{
	Use:   "myself",
	Short: "Show authenticated user",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newJiraClient()
		if err != nil {
			return err
		}
		user, err := client.GetMyself(cmd.Context())
		if err != nil {
			return err
		}
		return newFormatter(cmd).PrintUser(user)
	},
}

func init() {
	jiraCmd.AddCommand(jiraMyselfCmd)
}
