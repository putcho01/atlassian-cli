package cmd

import (
	"github.com/spf13/cobra"
)

var version = "0.2.2"

var rootCmd = &cobra.Command{
	Use:   "atlassian-cli",
	Short: "A CLI and MCP server for Atlassian Jira and Confluence",
	Long:  "A lightweight, native Go CLI for self-hosted Atlassian Jira and Confluence Server/Data Center instances.",
}

func init() {
	rootCmd.PersistentFlags().StringP("output", "o", "table", "Output format: table, json, markdown")
}

func Execute() error {
	return rootCmd.Execute()
}
