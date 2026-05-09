package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/putcho01/atlassian-cli/internal/config"
	"github.com/spf13/cobra"
)

var jiraIssueOpenCmd = &cobra.Command{
	Use:   "open <issue-key>",
	Short: "Open a Jira issue in the browser",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadJiraConfig()
		if err != nil {
			return err
		}
		url := strings.TrimRight(cfg.URL, "/") + "/browse/" + args[0]
		return openBrowser(url)
	},
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
