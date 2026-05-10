// atlassian-cli is a lightweight, native Go CLI and MCP server for
// self-hosted Atlassian Jira and Confluence Server/Data Center instances.
//
// # Installation
//
//	go install github.com/putcho01/atlassian-cli@latest
//
// # Configuration
//
// Set the following environment variables before use:
//
//	JIRA_URL        – Base URL of your Jira instance (e.g. https://jira.example.com)
//	JIRA_EMAIL      – Your Atlassian account email
//	JIRA_TOKEN      – API token or password
//	CONFLUENCE_URL  – Base URL of your Confluence instance
//	CONFLUENCE_EMAIL – Your Atlassian account email
//	CONFLUENCE_TOKEN – API token or password
//
// # Usage
//
//	atlassian-cli jira issue get PROJECT-123
//	atlassian-cli jira issue search "project = PROJECT AND status = 'In Progress'"
//	atlassian-cli confluence page get 12345678
//	atlassian-cli mcp-server
//
// # Output Formats
//
// Most commands support --output (-o) flag: table (default), json, markdown.
package main

import (
	"os"

	"github.com/putcho01/atlassian-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
