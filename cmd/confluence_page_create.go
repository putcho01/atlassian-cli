package cmd

import (
	"github.com/putcho01/atlassian-cli/internal/confluence"
	"github.com/spf13/cobra"
)

var confluencePageCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Confluence page",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newConfluenceClient()
		if err != nil {
			return err
		}

		space, _ := cmd.Flags().GetString("space")
		title, _ := cmd.Flags().GetString("title")
		body, _ := cmd.Flags().GetString("body")
		parent, _ := cmd.Flags().GetString("parent")

		page, err := client.CreatePage(cmd.Context(), &confluence.CreatePageInput{
			SpaceKey: space,
			Title:    title,
			Body:     body,
			ParentID: parent,
		})
		if err != nil {
			return err
		}
		return newFormatter(cmd).PrintPage(page)
	},
}

func init() {
	confluencePageCreateCmd.Flags().String("space", "", "Space key (required)")
	confluencePageCreateCmd.Flags().String("title", "", "Page title (required)")
	confluencePageCreateCmd.Flags().String("body", "", "Page body in Confluence storage format")
	confluencePageCreateCmd.Flags().String("parent", "", "Parent page ID")
	_ = confluencePageCreateCmd.MarkFlagRequired("space")
	_ = confluencePageCreateCmd.MarkFlagRequired("title")
	confluencePageCmd.AddCommand(confluencePageCreateCmd)
}
