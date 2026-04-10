package cmd

import (
	"github.com/spf13/cobra"
)

var confluenceCmd = &cobra.Command{
	Use:     "confluence",
	Aliases: []string{"conf"},
	Short:   "Confluence commands",
}

func init() {
	rootCmd.AddCommand(confluenceCmd)
}
