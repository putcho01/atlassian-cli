package mcptools

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/putcho01/atlassian-cli/internal/config"
	"github.com/putcho01/atlassian-cli/internal/htmlconv"
	"github.com/putcho01/atlassian-cli/internal/jira"
)

// --- Input types ---

type emptyInput struct{}

type getIssueInput struct {
	Key string `json:"key" jsonschema:"description=The Jira issue key (e.g. PROJ-123)"`
}

type searchIssuesInput struct {
	JQL        string `json:"jql" jsonschema:"description=JQL query string"`
	MaxResults int    `json:"max_results,omitempty" jsonschema:"description=Maximum number of results (default 50)"`
}

type createIssueInput struct {
	ProjectKey  string `json:"project_key" jsonschema:"description=Project key (e.g. PROJ)"`
	Summary     string `json:"summary" jsonschema:"description=Issue summary"`
	IssueType   string `json:"issue_type,omitempty" jsonschema:"description=Issue type (default Task)"`
	Description string `json:"description,omitempty" jsonschema:"description=Issue description"`
	Assignee    string `json:"assignee,omitempty" jsonschema:"description=Assignee username"`
	Priority    string `json:"priority,omitempty" jsonschema:"description=Priority name"`
}

type updateIssueInput struct {
	Key    string         `json:"key" jsonschema:"description=The Jira issue key"`
	Fields map[string]any `json:"fields" jsonschema:"description=Fields to update as key-value pairs"`
}

type deleteIssueInput struct {
	Key string `json:"key" jsonschema:"description=The Jira issue key to delete"`
}

type transitionIssueInput struct {
	Key    string `json:"key" jsonschema:"description=The Jira issue key"`
	Status string `json:"status" jsonschema:"description=Target status name"`
}

type getTransitionsInput struct {
	Key string `json:"key" jsonschema:"description=The Jira issue key"`
}

type getSubtasksInput struct {
	Key string `json:"key" jsonschema:"description=The parent Jira issue key"`
}

// --- Registration ---

func registerJiraTools(s *mcp.Server, enabledGroups []string) {
	maybeAdd(s, enabledGroups, "jira_get_myself", "Get the currently authenticated Jira user",
		func(ctx context.Context, req *mcp.CallToolRequest, input emptyInput) (*mcp.CallToolResult, any, error) {
			client, err := newJiraClient()
			if err != nil {
				return nil, nil, err
			}
			user, err := client.GetMyself(ctx)
			if err != nil {
				return nil, nil, err
			}
			return textResult(user)
		})

	maybeAdd(s, enabledGroups, "jira_get_issue", "Get a Jira issue by key, including summary, description, status, assignee, and other fields",
		func(ctx context.Context, req *mcp.CallToolRequest, input getIssueInput) (*mcp.CallToolResult, any, error) {
			client, err := newJiraClient()
			if err != nil {
				return nil, nil, err
			}
			issue, err := client.GetIssue(ctx, input.Key)
			if err != nil {
				return nil, nil, err
			}
			issue.Fields.Description = htmlconv.Convert(issue.Fields.Description)
			return textResult(issue)
		})

	maybeAdd(s, enabledGroups, "jira_search_issues", "Search Jira issues using JQL query language",
		func(ctx context.Context, req *mcp.CallToolRequest, input searchIssuesInput) (*mcp.CallToolResult, any, error) {
			client, err := newJiraClient()
			if err != nil {
				return nil, nil, err
			}
			result, err := client.SearchIssues(ctx, input.JQL, input.MaxResults)
			if err != nil {
				return nil, nil, err
			}
			for i := range result.Issues {
				result.Issues[i].Fields.Description = htmlconv.Convert(result.Issues[i].Fields.Description)
			}
			return textResult(result)
		})

	maybeAdd(s, enabledGroups, "jira_create_issue", "Create a new Jira issue",
		func(ctx context.Context, req *mcp.CallToolRequest, input createIssueInput) (*mcp.CallToolResult, any, error) {
			client, err := newJiraClient()
			if err != nil {
				return nil, nil, err
			}
			issueType := input.IssueType
			if issueType == "" {
				issueType = "Task"
			}
			createInput := &jira.CreateIssueInput{
				ProjectKey:  input.ProjectKey,
				Summary:     input.Summary,
				IssueType:   issueType,
				Description: input.Description,
				Assignee:    input.Assignee,
				Priority:    input.Priority,
			}
			issue, err := client.CreateIssue(ctx, createInput)
			if err != nil {
				return nil, nil, err
			}
			return textResult(issue)
		})

	maybeAdd(s, enabledGroups, "jira_update_issue", "Update fields on an existing Jira issue",
		func(ctx context.Context, req *mcp.CallToolRequest, input updateIssueInput) (*mcp.CallToolResult, any, error) {
			client, err := newJiraClient()
			if err != nil {
				return nil, nil, err
			}
			if err := client.UpdateIssue(ctx, input.Key, input.Fields); err != nil {
				return nil, nil, err
			}
			return textResult(map[string]string{"message": "Issue " + input.Key + " updated successfully"})
		})

	maybeAdd(s, enabledGroups, "jira_delete_issue", "Delete a Jira issue",
		func(ctx context.Context, req *mcp.CallToolRequest, input deleteIssueInput) (*mcp.CallToolResult, any, error) {
			client, err := newJiraClient()
			if err != nil {
				return nil, nil, err
			}
			if err := client.DeleteIssue(ctx, input.Key); err != nil {
				return nil, nil, err
			}
			return textResult(map[string]string{"message": "Issue " + input.Key + " deleted"})
		})

	maybeAdd(s, enabledGroups, "jira_get_transitions", "Get available transitions for a Jira issue",
		func(ctx context.Context, req *mcp.CallToolRequest, input getTransitionsInput) (*mcp.CallToolResult, any, error) {
			client, err := newJiraClient()
			if err != nil {
				return nil, nil, err
			}
			transitions, err := client.GetTransitions(ctx, input.Key)
			if err != nil {
				return nil, nil, err
			}
			return textResult(transitions)
		})

	maybeAdd(s, enabledGroups, "jira_transition_issue", "Transition a Jira issue to a new status",
		func(ctx context.Context, req *mcp.CallToolRequest, input transitionIssueInput) (*mcp.CallToolResult, any, error) {
			client, err := newJiraClient()
			if err != nil {
				return nil, nil, err
			}
			if err := client.TransitionIssue(ctx, input.Key, input.Status); err != nil {
				return nil, nil, err
			}
			return textResult(map[string]string{"message": "Issue " + input.Key + " transitioned to " + input.Status})
		})

	maybeAdd(s, enabledGroups, "jira_get_subtasks", "Get subtasks of a Jira issue",
		func(ctx context.Context, req *mcp.CallToolRequest, input getSubtasksInput) (*mcp.CallToolResult, any, error) {
			client, err := newJiraClient()
			if err != nil {
				return nil, nil, err
			}
			subtasks, err := client.GetSubtasks(ctx, input.Key)
			if err != nil {
				return nil, nil, err
			}
			return textResult(subtasks)
		})
}

func newJiraClient() (*jira.Client, error) {
	cfg, err := config.LoadJiraConfig()
	if err != nil {
		return nil, err
	}
	return jira.NewClient(cfg), nil
}

func textResult(v any) (*mcp.CallToolResult, any, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil, nil
}
