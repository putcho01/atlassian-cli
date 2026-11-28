package mcptools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/putcho01/atlassian-cli/internal/jira"
)

func TestJiraToolsKeepTextVerbatim(t *testing.T) {
	const text = "use List<String> here"
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		issue := jira.Issue{Key: "P-1", Fields: jira.IssueFields{Description: text}}
		var body any = issue
		if r.URL.Path == "/rest/api/2/search" {
			body = jira.SearchResult{Issues: []jira.Issue{issue}}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(body)
	}))
	defer api.Close()

	t.Setenv("JIRA_URL", api.URL)
	t.Setenv("JIRA_EMAIL", "")
	t.Setenv("JIRA_API_TOKEN", "test-token")

	ctx := context.Background()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	if _, err := NewServer("test", nil).Connect(ctx, serverTransport, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "test"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer session.Close()

	for name, args := range map[string]map[string]any{
		"jira_get_issue":     {"key": "P-1"},
		"jira_search_issues": {"jql": "project = P"},
	} {
		res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: args})
		if err != nil {
			t.Fatalf("call %s: %v", name, err)
		}
		var got struct {
			Fields jira.IssueFields `json:"fields"`
			Issues []jira.Issue     `json:"issues"`
		}
		if err := json.Unmarshal([]byte(res.Content[0].(*mcp.TextContent).Text), &got); err != nil {
			t.Fatalf("decode %s result: %v", name, err)
		}
		description := got.Fields.Description
		if len(got.Issues) > 0 {
			description = got.Issues[0].Fields.Description
		}
		if description != text {
			t.Errorf("%s description = %q, want %q", name, description, text)
		}
	}
}

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
