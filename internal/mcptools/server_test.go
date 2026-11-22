package mcptools

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestNewServerRegistersAllTools(t *testing.T) {
	ctx := context.Background()
	server := NewServer("test", nil)

	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, serverTransport, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "test"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer session.Close()

	res, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}

	want := 0
	for _, tools := range toolGroups {
		want += len(tools)
	}
	if len(res.Tools) != want {
		t.Errorf("got %d tools, want %d", len(res.Tools), want)
	}

	for _, tool := range res.Tools {
		schema, ok := tool.InputSchema.(map[string]any)
		if !ok {
			t.Fatalf("tool %s: unexpected input schema type %T", tool.Name, tool.InputSchema)
		}
		props, _ := schema["properties"].(map[string]any)
		for name, p := range props {
			desc, _ := p.(map[string]any)["description"].(string)
			if desc == "" || strings.HasPrefix(desc, "description=") {
				t.Errorf("tool %s property %s: bad description %q", tool.Name, name, desc)
			}
		}
	}
}
