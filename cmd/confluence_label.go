package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var confluenceLabelCmd = &cobra.Command{
	Use:   "label",
	Short: "Confluence label commands",
}

var confluenceLabelListCmd = &cobra.Command{
	Use:   "list <page-id>",
	Short: "List labels on a Confluence page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newConfluenceClient()
		if err != nil {
			return err
		}
		labels, err := client.GetLabels(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return newFormatter(cmd).PrintLabels(labels)
	},
}

var confluenceLabelAddCmd = &cobra.Command{
	Use:   "add <page-id> <label1,label2,...>",
	Short: "Add labels to a Confluence page",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newConfluenceClient()
		if err != nil {
			return err
		}
		labels := strings.Split(args[1], ",")
		if err := client.AddLabels(cmd.Context(), args[0], labels); err != nil {
			return err
		}
		return newFormatter(cmd).PrintMessage(fmt.Sprintf("Labels added: %s", args[1]))
	},
}

var confluenceLabelRemoveCmd = &cobra.Command{
	Use:   "remove <page-id> <label>",
	Short: "Remove a label from a Confluence page",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newConfluenceClient()
		if err != nil {
			return err
		}
		if err := client.RemoveLabel(cmd.Context(), args[0], args[1]); err != nil {
			return err
		}
		return newFormatter(cmd).PrintMessage(fmt.Sprintf("Label %q removed", args[1]))
	},
}

func init() {
	confluenceLabelCmd.AddCommand(confluenceLabelListCmd)
	confluenceLabelCmd.AddCommand(confluenceLabelAddCmd)
	confluenceLabelCmd.AddCommand(confluenceLabelRemoveCmd)
	confluenceCmd.AddCommand(confluenceLabelCmd)
}
