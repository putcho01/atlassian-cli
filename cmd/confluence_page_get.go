package cmd

import (
	"github.com/spf13/cobra"
)

var confluencePageGetCmd = &cobra.Command{
	Use:   "get <page-id>",
	Short: "Get a Confluence page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newConfluenceClient()
		if err != nil {
			return err
		}
		page, err := client.GetPage(cmd.Context(), args[0], nil)
		if err != nil {
			return err
		}
		return newFormatter(cmd).PrintPage(page)
	},
}

func init() {
	confluencePageCmd.AddCommand(confluencePageGetCmd)
}
