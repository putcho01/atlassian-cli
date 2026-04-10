package mcptools

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/putcho01/atlassian-cli/internal/config"
	"github.com/putcho01/atlassian-cli/internal/confluence"
	"github.com/putcho01/atlassian-cli/internal/htmlconv"
)

// --- Input types ---

type getPageInput struct {
	ID     string `json:"id" jsonschema:"description=The Confluence page ID"`
	Expand string `json:"expand,omitempty" jsonschema:"description=Comma-separated list of properties to expand (default: body.storage,version)"`
}

type listLabelsInput struct {
	PageID string `json:"page_id" jsonschema:"description=The Confluence page ID"`
}

type addLabelsInput struct {
	PageID string   `json:"page_id" jsonschema:"description=The Confluence page ID"`
	Labels []string `json:"labels" jsonschema:"description=Labels to add"`
}

type removeLabelInput struct {
	PageID string `json:"page_id" jsonschema:"description=The Confluence page ID"`
	Label  string `json:"label" jsonschema:"description=Label name to remove"`
}

type listRestrictionsInput struct {
	PageID string `json:"page_id" jsonschema:"description=The Confluence page ID"`
}

type restrictionInput struct {
	PageID    string `json:"page_id" jsonschema:"description=The Confluence page ID"`
	Operation string `json:"operation" jsonschema:"description=Operation type: read or update"`
	Type      string `json:"type" jsonschema:"description=Restriction type: user or group"`
	Name      string `json:"name" jsonschema:"description=Username or group name"`
}

// --- Registration ---

func registerConfluenceTools(s *mcp.Server, enabledGroups []string) {
	maybeAdd(s, enabledGroups, "confluence_get_page", "Get a Confluence page by ID, including its content",
		func(ctx context.Context, req *mcp.CallToolRequest, input getPageInput) (*mcp.CallToolResult, any, error) {
			client, err := newConfluenceClient()
			if err != nil {
				return nil, nil, err
			}
			var expand []string
			if input.Expand != "" {
				expand = strings.Split(input.Expand, ",")
			}
			page, err := client.GetPage(ctx, input.ID, expand)
			if err != nil {
				return nil, nil, err
			}
			if page.Body.Storage != nil {
				page.Body.Storage.Value = htmlconv.Convert(page.Body.Storage.Value)
			}
			return textResult(page)
		})

	maybeAdd(s, enabledGroups, "confluence_list_labels", "List labels on a Confluence page",
		func(ctx context.Context, req *mcp.CallToolRequest, input listLabelsInput) (*mcp.CallToolResult, any, error) {
			client, err := newConfluenceClient()
			if err != nil {
				return nil, nil, err
			}
			labels, err := client.GetLabels(ctx, input.PageID)
			if err != nil {
				return nil, nil, err
			}
			return textResult(labels)
		})

	maybeAdd(s, enabledGroups, "confluence_add_labels", "Add labels to a Confluence page",
		func(ctx context.Context, req *mcp.CallToolRequest, input addLabelsInput) (*mcp.CallToolResult, any, error) {
			client, err := newConfluenceClient()
			if err != nil {
				return nil, nil, err
			}
			if err := client.AddLabels(ctx, input.PageID, input.Labels); err != nil {
				return nil, nil, err
			}
			return textResult(map[string]string{"message": "Labels added successfully"})
		})

	maybeAdd(s, enabledGroups, "confluence_remove_label", "Remove a label from a Confluence page",
		func(ctx context.Context, req *mcp.CallToolRequest, input removeLabelInput) (*mcp.CallToolResult, any, error) {
			client, err := newConfluenceClient()
			if err != nil {
				return nil, nil, err
			}
			if err := client.RemoveLabel(ctx, input.PageID, input.Label); err != nil {
				return nil, nil, err
			}
			return textResult(map[string]string{"message": "Label removed successfully"})
		})

	maybeAdd(s, enabledGroups, "confluence_list_restrictions", "List restrictions on a Confluence page",
		func(ctx context.Context, req *mcp.CallToolRequest, input listRestrictionsInput) (*mcp.CallToolResult, any, error) {
			client, err := newConfluenceClient()
			if err != nil {
				return nil, nil, err
			}
			restrictions, err := client.GetRestrictions(ctx, input.PageID)
			if err != nil {
				return nil, nil, err
			}
			return textResult(restrictions)
		})

	maybeAdd(s, enabledGroups, "confluence_add_restriction", "Add a restriction to a Confluence page",
		func(ctx context.Context, req *mcp.CallToolRequest, input restrictionInput) (*mcp.CallToolResult, any, error) {
			client, err := newConfluenceClient()
			if err != nil {
				return nil, nil, err
			}
			ri := &confluence.RestrictionInput{
				Operation: input.Operation,
				Type:      input.Type,
				Name:      input.Name,
			}
			if err := client.AddRestriction(ctx, input.PageID, ri); err != nil {
				return nil, nil, err
			}
			return textResult(map[string]string{"message": "Restriction added successfully"})
		})

	maybeAdd(s, enabledGroups, "confluence_remove_restriction", "Remove a restriction from a Confluence page",
		func(ctx context.Context, req *mcp.CallToolRequest, input restrictionInput) (*mcp.CallToolResult, any, error) {
			client, err := newConfluenceClient()
			if err != nil {
				return nil, nil, err
			}
			ri := &confluence.RestrictionInput{
				Operation: input.Operation,
				Type:      input.Type,
				Name:      input.Name,
			}
			if err := client.RemoveRestriction(ctx, input.PageID, ri); err != nil {
				return nil, nil, err
			}
			return textResult(map[string]string{"message": "Restriction removed successfully"})
		})
}

func newConfluenceClient() (*confluence.Client, error) {
	cfg, err := config.LoadConfluenceConfig()
	if err != nil {
		return nil, err
	}
	return confluence.NewClient(cfg), nil
}
