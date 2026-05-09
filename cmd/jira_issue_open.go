package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var jiraIssueOpenCmd = &cobra.Command{
	Use:   "open <issue-key>",
	Short: "Open a Jira issue in the browser",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newJiraClient()
		if err != nil {
			return err
		}
		return openBrowser(issueURL(client.BaseURL(), args[0]))
	},
}

func issueURL(baseURL, key string) string {
	return baseURL + "/browse/" + key
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}
	return nil
}

func init() {
	jiraIssueCmd.AddCommand(jiraIssueOpenCmd)
}
