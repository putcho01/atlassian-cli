package cmd

import (
	"fmt"

	"github.com/putcho01/atlassian-cli/internal/confluence"
	"github.com/spf13/cobra"
)

var confluenceRestrictionCmd = &cobra.Command{
	Use:   "restriction",
	Short: "Confluence page restriction commands",
}

var confluenceRestrictionListCmd = &cobra.Command{
	Use:   "list <page-id>",
	Short: "List restrictions on a Confluence page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newConfluenceClient()
		if err != nil {
			return err
		}
		restrictions, err := client.GetRestrictions(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return newFormatter(cmd).PrintJSON(restrictions)
	},
}

var confluenceRestrictionAddCmd = &cobra.Command{
	Use:   "add <page-id>",
	Short: "Add a restriction to a Confluence page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newConfluenceClient()
		if err != nil {
			return err
		}
		if err := client.AddRestriction(cmd.Context(), args[0], parseRestrictionInput(cmd)); err != nil {
			return err
		}
		return newFormatter(cmd).PrintMessage(fmt.Sprintf("Restriction added to page %s", args[0]))
	},
}

var confluenceRestrictionRemoveCmd = &cobra.Command{
	Use:   "remove <page-id>",
	Short: "Remove a restriction from a Confluence page",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newConfluenceClient()
		if err != nil {
			return err
		}
		if err := client.RemoveRestriction(cmd.Context(), args[0], parseRestrictionInput(cmd)); err != nil {
			return err
		}
		return newFormatter(cmd).PrintMessage(fmt.Sprintf("Restriction removed from page %s", args[0]))
	},
}

func parseRestrictionInput(cmd *cobra.Command) *confluence.RestrictionInput {
	operation, _ := cmd.Flags().GetString("operation")
	restrictionType, _ := cmd.Flags().GetString("type")
	name, _ := cmd.Flags().GetString("name")
	return &confluence.RestrictionInput{
		Operation: operation,
		Type:      restrictionType,
		Name:      name,
	}
}

func init() {
	for _, cmd := range []*cobra.Command{confluenceRestrictionAddCmd, confluenceRestrictionRemoveCmd} {
		cmd.Flags().String("operation", "update", "Operation: read or update")
		cmd.Flags().String("type", "user", "Restriction type: user or group")
		cmd.Flags().String("name", "", "Username or group name")
		_ = cmd.MarkFlagRequired("name")
	}
	confluenceRestrictionCmd.AddCommand(confluenceRestrictionListCmd)
	confluenceRestrictionCmd.AddCommand(confluenceRestrictionAddCmd)
	confluenceRestrictionCmd.AddCommand(confluenceRestrictionRemoveCmd)
	confluenceCmd.AddCommand(confluenceRestrictionCmd)
}
