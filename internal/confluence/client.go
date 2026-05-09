package confluence

import (
	"context"
	"fmt"
	"io"
	"net/url"
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

func (c *Client) GetPage(ctx context.Context, id string, expand []string) (*Page, error) {
	q := url.Values{}
	if len(expand) > 0 {
		q.Set("expand", strings.Join(expand, ","))
	} else {
		q.Set("expand", "body.storage,version")
	}

	var page Page
	err := c.http.Do(ctx, "GET", contentPath(id), q, nil, &page)
	return &page, err
}

func (c *Client) GetLabels(ctx context.Context, pageID string) ([]Label, error) {
	var resp LabelResponse
	err := c.http.Do(ctx, "GET", contentPath(pageID)+"/label", nil, nil, &resp)
	return resp.Results, err
}

func (c *Client) AddLabels(ctx context.Context, pageID string, labels []string) error {
	body := make([]map[string]string, len(labels))
	for i, l := range labels {
		body[i] = map[string]string{"prefix": "global", "name": l}
	}
	return c.http.Do(ctx, "POST", contentPath(pageID)+"/label", nil, body, nil)
}

func (c *Client) RemoveLabel(ctx context.Context, pageID, label string) error {
	q := url.Values{}
	q.Set("name", label)
	return c.http.Do(ctx, "DELETE", contentPath(pageID)+"/label", q, nil, nil)
}

func (c *Client) GetRestrictions(ctx context.Context, pageID string) ([]Restriction, error) {
	var result []Restriction
	err := c.http.Do(ctx, "GET", contentPath(pageID)+"/restriction", nil, nil, &result)
	return result, err
}

func (c *Client) AddRestriction(ctx context.Context, pageID string, input *RestrictionInput) error {
	body := buildRestrictionBody(input)
	return c.http.Do(ctx, "PUT", contentPath(pageID)+"/restriction", nil, body, nil)
}

func (c *Client) RemoveRestriction(ctx context.Context, pageID string, input *RestrictionInput) error {
	q := url.Values{}
	var path string
	if input.Type == "user" {
		path = contentPath(pageID) + "/restriction/byOperation/" + url.PathEscape(input.Operation) + "/user"
		q.Set("accountId", input.Name)
	} else {
		path = contentPath(pageID) + "/restriction/byOperation/" + url.PathEscape(input.Operation) + "/group/" + url.PathEscape(input.Name)
	}
	return c.http.Do(ctx, "DELETE", path, q, nil, nil)
}

func (c *Client) CreatePage(ctx context.Context, input *CreatePageInput) (*Page, error) {
	body := map[string]any{
		"type":  "page",
		"title": input.Title,
		"space": map[string]string{"key": input.SpaceKey},
		"body":  storageBody(input.Body),
	}
	if input.ParentID != "" {
		body["ancestors"] = []map[string]string{{"id": input.ParentID}}
	}
	var page Page
	if err := c.http.Do(ctx, "POST", "/rest/api/content", nil, body, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

func (c *Client) UpdatePage(ctx context.Context, id string, input *UpdatePageInput) (*Page, error) {
	version := input.Version
	if version == 0 {
		current, err := c.GetPage(ctx, id, []string{"version"})
		if err != nil {
			return nil, fmt.Errorf("fetch current version: %w", err)
		}
		version = current.Version.Number + 1
	}
	body := map[string]any{
		"type":  "page",
		"title": input.Title,
		"version": map[string]any{
			"number":  version,
			"message": input.Message,
		},
		"body": storageBody(input.Body),
	}
	var page Page
	if err := c.http.Do(ctx, "PUT", contentPath(id), nil, body, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

func (c *Client) GetAttachments(ctx context.Context, pageID string) ([]Attachment, error) {
	var resp AttachmentResponse
	err := c.http.Do(ctx, "GET", contentPath(pageID)+"/child/attachment", nil, nil, &resp)
	return resp.Results, err
}

func (c *Client) DownloadAttachment(ctx context.Context, downloadPath string) (io.ReadCloser, error) {
	resp, err := c.http.Get(ctx, downloadPath, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}
	return resp.Body, nil
}

func contentPath(id string) string {
	return "/rest/api/content/" + url.PathEscape(id)
}

func storageBody(value string) map[string]any {
	return map[string]any{
		"storage": map[string]string{
			"value":          value,
			"representation": "storage",
		},
	}
}

func buildRestrictionBody(input *RestrictionInput) []map[string]any {
	subject := map[string]any{}
	if input.Type == "user" {
		subject["user"] = map[string]any{
			"type":      "known",
			"accountId": input.Name,
		}
	} else {
		subject["group"] = map[string]any{
			"type": "group",
			"name": input.Name,
		}
	}
	return []map[string]any{
		{
			"operation":    input.Operation,
			"restrictions": subject,
		},
	}
}
