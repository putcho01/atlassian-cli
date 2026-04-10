package cmd

import (
	"os"

	"github.com/putcho01/atlassian-cli/internal/config"
	"github.com/putcho01/atlassian-cli/internal/confluence"
	"github.com/putcho01/atlassian-cli/internal/formatter"
	"github.com/putcho01/atlassian-cli/internal/jira"
	"github.com/spf13/cobra"
)

func newJiraClient() (*jira.Client, error) {
	cfg, err := config.LoadJiraConfig()
	if err != nil {
		return nil, err
	}
	return jira.NewClient(cfg), nil
}

func newConfluenceClient() (*confluence.Client, error) {
	cfg, err := config.LoadConfluenceConfig()
	if err != nil {
		return nil, err
	}
	return confluence.NewClient(cfg), nil
}

func newFormatter(cmd *cobra.Command) *formatter.Formatter {
	format, _ := cmd.Flags().GetString("output")
	return formatter.New(formatter.ParseFormat(format), os.Stdout)
}
