package jira

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/putcho01/atlassian-cli/internal/config"
	"github.com/putcho01/atlassian-cli/internal/httpclient"
)

type Client struct {
	http *httpclient.Client
}

func NewClient(cfg *config.AtlassianConfig) *Client {
	return &Client{http: httpclient.NewFromConfig(cfg.URL, cfg.Email, cfg.Token)}
}

func (c *Client) GetMyself(ctx context.Context) (*User, error) {
	var user User
	err := c.http.Do(ctx, "GET", "/rest/api/2/myself", nil, nil, &user)
	return &user, err
}

func (c *Client) GetIssue(ctx context.Context, key string) (*Issue, error) {
	var issue Issue
	err := c.http.Do(ctx, "GET", "/rest/api/2/issue/"+url.PathEscape(key), nil, nil, &issue)
	return &issue, err
}

func (c *Client) SearchIssues(ctx context.Context, jql string, maxResults int) (*SearchResult, error) {
	if maxResults <= 0 {
		maxResults = 50
	}
	q := url.Values{}
	q.Set("jql", jql)
	q.Set("maxResults", strconv.Itoa(maxResults))

	var result SearchResult
	err := c.http.Do(ctx, "GET", "/rest/api/2/search", q, nil, &result)
	return &result, err
}

func (c *Client) CreateIssue(ctx context.Context, input *CreateIssueInput) (*Issue, error) {
	body := map[string]any{
		"fields": buildCreateFields(input),
	}

	var issue Issue
	err := c.http.Do(ctx, "POST", "/rest/api/2/issue", nil, body, &issue)
	return &issue, err
}

func (c *Client) UpdateIssue(ctx context.Context, key string, fields map[string]any) error {
	body := map[string]any{"fields": fields}
	return c.http.Do(ctx, "PUT", "/rest/api/2/issue/"+url.PathEscape(key), nil, body, nil)
}

func (c *Client) DeleteIssue(ctx context.Context, key string) error {
	return c.http.Do(ctx, "DELETE", "/rest/api/2/issue/"+url.PathEscape(key), nil, nil, nil)
}

func (c *Client) GetTransitions(ctx context.Context, key string) ([]Transition, error) {
	var resp transitionsResponse
	err := c.http.Do(ctx, "GET", "/rest/api/2/issue/"+url.PathEscape(key)+"/transitions", nil, nil, &resp)
	return resp.Transitions, err
}

func (c *Client) TransitionIssue(ctx context.Context, key, statusName string) error {
	transitions, err := c.GetTransitions(ctx, key)
	if err != nil {
		return fmt.Errorf("get transitions: %w", err)
	}

	var transitionID string
	for _, t := range transitions {
		if strings.EqualFold(t.Name, statusName) || (t.To != nil && strings.EqualFold(t.To.Name, statusName)) {
			transitionID = t.ID
			break
		}
	}
	if transitionID == "" {
		available := make([]string, len(transitions))
		for i, t := range transitions {
			available[i] = t.Name
		}
		return fmt.Errorf("transition %q not found, available: %s", statusName, strings.Join(available, ", "))
	}

	body := map[string]any{
		"transition": map[string]string{"id": transitionID},
	}
	return c.http.Do(ctx, "POST", "/rest/api/2/issue/"+url.PathEscape(key)+"/transitions", nil, body, nil)
}

func (c *Client) ListComments(ctx context.Context, key string) ([]Comment, error) {
	var resp commentsResponse
	err := c.http.Do(ctx, "GET", "/rest/api/2/issue/"+url.PathEscape(key)+"/comment", nil, nil, &resp)
	return resp.Comments, err
}

func (c *Client) AddComment(ctx context.Context, key, body string) (*Comment, error) {
	payload := map[string]any{"body": body}
	var comment Comment
	if err := c.http.Do(ctx, "POST", "/rest/api/2/issue/"+url.PathEscape(key)+"/comment", nil, payload, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

func (c *Client) GetSubtasks(ctx context.Context, key string) ([]Issue, error) {
	issue, err := c.GetIssue(ctx, key)
	if err != nil {
		return nil, err
	}
	return issue.Fields.Subtasks, nil
}

func buildCreateFields(input *CreateIssueInput) map[string]any {
	fields := map[string]any{
		"project":   map[string]string{"key": input.ProjectKey},
		"summary":   input.Summary,
		"issuetype": map[string]string{"name": input.IssueType},
	}
	if input.Description != "" {
		fields["description"] = input.Description
	}
	if input.Assignee != "" {
		fields["assignee"] = map[string]string{"accountId": input.Assignee}
	}
	if input.Priority != "" {
		fields["priority"] = map[string]string{"name": input.Priority}
	}
	if len(input.Labels) > 0 {
		fields["labels"] = input.Labels
	}
	for k, v := range input.CustomFields {
		fields[k] = v
	}
	return fields
}
