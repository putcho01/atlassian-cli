package cmd

import (
	"github.com/spf13/cobra"
)

var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "Jira commands",
}

func init() {
	rootCmd.AddCommand(jiraCmd)
}
