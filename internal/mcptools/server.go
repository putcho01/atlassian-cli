package mcptools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NewServer creates a new MCP server with Atlassian tools registered.
// enabledGroups filters which tool groups are enabled. Empty means all.
func NewServer(version string, enabledGroups []string) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "atlassian-cli",
		Version: version,
	}, nil)

	registerJiraTools(server, enabledGroups)
	registerConfluenceTools(server, enabledGroups)

	return server
}

// Run starts the MCP server with stdio transport.
func Run(ctx context.Context, version string, enabledGroups []string) error {
	server := NewServer(version, enabledGroups)
	return server.Run(ctx, &mcp.StdioTransport{})
}

// maybeAdd registers a tool if it's in the enabled groups.
func maybeAdd[In, Out any](s *mcp.Server, enabledGroups []string, name, description string, handler mcp.ToolHandlerFor[In, Out]) {
	if !isToolEnabled(name, enabledGroups) {
		return
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        name,
		Description: description,
	}, handler)
}
