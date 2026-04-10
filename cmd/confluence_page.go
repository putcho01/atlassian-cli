package cmd

import (
	"github.com/spf13/cobra"
)

var confluencePageCmd = &cobra.Command{
	Use:   "page",
	Short: "Confluence page commands",
}

func init() {
	confluenceCmd.AddCommand(confluencePageCmd)
}
