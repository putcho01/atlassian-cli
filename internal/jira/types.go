package jira

// User represents a Jira user account.
type User struct {
	Self         string `json:"self"`
	AccountID    string `json:"accountId"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
	TimeZone     string `json:"timeZone"`
}

// Issue represents a Jira issue with its key and fields.
type Issue struct {
	ID     string      `json:"id"`
	Key    string      `json:"key"`
	Self   string      `json:"self"`
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Summary     string     `json:"summary"`
	Description string     `json:"description"`
	Status      *Status    `json:"status"`
	IssueType   *IssueType `json:"issuetype"`
	Priority    *Priority  `json:"priority"`
	Assignee    *User      `json:"assignee"`
	Reporter    *User      `json:"reporter"`
	Project     *Project   `json:"project"`
	Created     string     `json:"created"`
	Updated     string     `json:"updated"`
	Labels      []string   `json:"labels"`
	Subtasks    []Issue    `json:"subtasks"`
}

type Status struct {
	Self        string `json:"self"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type IssueType struct {
	Self        string `json:"self"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Subtask     bool   `json:"subtask"`
}

type Priority struct {
	Self string `json:"self"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Project struct {
	Self string `json:"self"`
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Transition struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	To   *Status `json:"to"`
}

// SearchResult holds a paginated list of issues returned by a JQL query.
type SearchResult struct {
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

// CreateIssueInput holds the parameters for creating a new Jira issue.
type CreateIssueInput struct {
	ProjectKey  string            `json:"projectKey"`
	Summary     string            `json:"summary"`
	IssueType   string            `json:"issueType"`
	Description string            `json:"description,omitempty"`
	Assignee    string            `json:"assignee,omitempty"`
	Priority    string            `json:"priority,omitempty"`
	Labels      []string          `json:"labels,omitempty"`
	CustomFields map[string]any   `json:"customFields,omitempty"`
}

type transitionsResponse struct {
	Transitions []Transition `json:"transitions"`
}

// Comment represents a comment on a Jira issue.
type Comment struct {
	ID      string `json:"id"`
	Author  *User  `json:"author"`
	Body    string `json:"body"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

type commentsResponse struct {
	Comments   []Comment `json:"comments"`
	Total      int       `json:"total"`
	StartAt    int       `json:"startAt"`
	MaxResults int       `json:"maxResults"`
}
