// Package cmd implements the atlassian-cli command-line interface.
//
// The command tree is organised as follows:
//
//	atlassian-cli
//	├── jira
//	│   ├── issue
//	│   │   ├── get        – fetch a single issue by key
//	│   │   ├── search     – search issues with JQL
//	│   │   ├── create     – create a new issue
//	│   │   ├── update     – update issue fields
//	│   │   ├── delete     – delete an issue
//	│   │   ├── transition – move an issue to a new status
//	│   │   ├── subtasks   – list subtasks of an issue
//	│   │   ├── open       – open an issue in the browser
//	│   │   └── comment
//	│   │       ├── list   – list comments on an issue
//	│   │       └── add    – add a comment to an issue
//	│   └── myself         – show the authenticated user
//	├── confluence (alias: conf)
//	│   ├── page
//	│   │   ├── get        – fetch a page by ID
//	│   │   ├── create     – create a new page
//	│   │   └── update     – update an existing page
//	│   ├── label
//	│   │   ├── list       – list labels on a page
//	│   │   ├── add        – add labels to a page
//	│   │   └── remove     – remove a label from a page
//	│   └── restriction
//	│       ├── list       – list restrictions on a page
//	│       ├── add        – add a read/update restriction
//	│       └── remove     – remove a restriction
//	├── mcp-server         – start an MCP server over stdio
//	└── version            – print version information
package cmd

import (
	"github.com/spf13/cobra"
)

var version = "0.2.7"

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
