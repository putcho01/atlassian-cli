package confluence

type Page struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`
	Status  string   `json:"status"`
	Title   string   `json:"title"`
	Body    PageBody `json:"body"`
	Version Version  `json:"version"`
	Links   Links    `json:"_links"`
}

type PageBody struct {
	Storage *BodyContent `json:"storage,omitempty"`
	View    *BodyContent `json:"view,omitempty"`
}

type BodyContent struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

type Version struct {
	Number  int    `json:"number"`
	When    string `json:"when"`
	Message string `json:"message"`
	By      *User  `json:"by"`
}

type User struct {
	Type        string `json:"type"`
	AccountID   string `json:"accountId"`
	DisplayName string `json:"displayName"`
}

type Links struct {
	WebUI string `json:"webui"`
	Self  string `json:"self"`
	Base  string `json:"base"`
}

type Label struct {
	Prefix string `json:"prefix"`
	Name   string `json:"name"`
	ID     string `json:"id"`
}

type LabelResponse struct {
	Results []Label `json:"results"`
	Size    int     `json:"size"`
}

type Restriction struct {
	Operation string              `json:"operation"`
	User      *RestrictionSubject `json:"restrictions,omitempty"`
}

type RestrictionSubject struct {
	User  *RestrictionUsers  `json:"user,omitempty"`
	Group *RestrictionGroups `json:"group,omitempty"`
}

type RestrictionUsers struct {
	Results []User `json:"results"`
	Size    int    `json:"size"`
}

type RestrictionGroups struct {
	Results []Group `json:"results"`
	Size    int     `json:"size"`
}

type Group struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type RestrictionInput struct {
	Operation string // "update" or "read"
	Type      string // "user" or "group"
	Name      string // username or group name
}

type Attachment struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Status    string            `json:"status"`
	Title     string            `json:"title"`
	MediaType string            `json:"mediaType"`
	FileSize  int64             `json:"fileSize"`
	Links     map[string]string `json:"_links"`
}

type AttachmentResponse struct {
	Results []Attachment `json:"results"`
	Size    int          `json:"size"`
}
