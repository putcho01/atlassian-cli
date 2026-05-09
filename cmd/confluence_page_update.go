package cmd

import (
	"github.com/putcho01/atlassian-cli/internal/confluence"
	"github.com/spf13/cobra"
)

var confluencePageUpdateCmd = &cobra.Command{
	Use:   "update <page-id>",
	Short: "Update a Confluence page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newConfluenceClient()
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		body, _ := cmd.Flags().GetString("body")
		version, _ := cmd.Flags().GetInt("version")
		message, _ := cmd.Flags().GetString("message")

		page, err := client.UpdatePage(cmd.Context(), args[0], &confluence.UpdatePageInput{
			Title:   title,
			Body:    body,
			Version: version,
			Message: message,
		})
		if err != nil {
			return err
		}
		return newFormatter(cmd).PrintPage(page)
	},
}

func init() {
	confluencePageUpdateCmd.Flags().String("title", "", "Page title")
	confluencePageUpdateCmd.Flags().String("body", "", "Page body in Confluence storage format")
	confluencePageUpdateCmd.Flags().Int("version", 0, "Version number to set (default: auto-fetch current version + 1)")
	confluencePageUpdateCmd.Flags().String("message", "", "Version message")
	_ = confluencePageUpdateCmd.MarkFlagRequired("title")
	_ = confluencePageUpdateCmd.MarkFlagRequired("body")
	confluencePageCmd.AddCommand(confluencePageUpdateCmd)
}
