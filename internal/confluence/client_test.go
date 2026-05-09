package confluence

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/putcho01/atlassian-cli/internal/httpclient"
)

func newTestConfluenceClient(baseURL string) *Client {
	return &Client{http: httpclient.NewFromConfig(baseURL, "test@test.com", "token")}
}

func TestCreatePage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   *CreatePageInput
		want    *Page
		wantErr bool
	}{
		{
			name: "success without parent",
			input: &CreatePageInput{
				SpaceKey: "PROJ",
				Title:    "New Page",
				Body:     "<p>Hello</p>",
			},
			want: &Page{ID: "123", Title: "New Page", Status: "current"},
		},
		{
			name: "success with parent",
			input: &CreatePageInput{
				SpaceKey: "PROJ",
				Title:    "Child Page",
				Body:     "<p>Child</p>",
				ParentID: "99",
			},
			want: &Page{ID: "124", Title: "Child Page", Status: "current"},
		},
		{
			name:    "server error",
			input:   &CreatePageInput{SpaceKey: "ERR", Title: "fail", Body: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mux := http.NewServeMux()
			mux.HandleFunc("/rest/api/content", func(w http.ResponseWriter, r *http.Request) {
				var req map[string]any
				json.NewDecoder(r.Body).Decode(&req)
				if space, _ := req["space"].(map[string]any); space != nil {
					if space["key"] == "ERR" {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tt.want)
			})
			srv := httptest.NewServer(mux)
			defer srv.Close()

			got, err := newTestConfluenceClient(srv.URL).CreatePage(context.Background(), tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("CreatePage() diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdatePage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		pageID  string
		input   *UpdatePageInput
		want    *Page
		wantErr bool
	}{
		{
			name:   "success",
			pageID: "123",
			input:  &UpdatePageInput{Title: "Updated", Body: "<p>new</p>", Version: 2},
			want:   &Page{ID: "123", Title: "Updated", Status: "current"},
		},
		{
			name:    "server error",
			pageID:  "999",
			input:   &UpdatePageInput{Title: "fail", Body: "", Version: 1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mux := http.NewServeMux()
			mux.HandleFunc("/rest/api/content/123", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tt.want)
			})
			mux.HandleFunc("/rest/api/content/999", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			})
			srv := httptest.NewServer(mux)
			defer srv.Close()

			got, err := newTestConfluenceClient(srv.URL).UpdatePage(context.Background(), tt.pageID, tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("UpdatePage() diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdatePageAutoVersion(t *testing.T) {
	t.Parallel()
	current := &Page{ID: "123", Title: "Old Title", Status: "current", Version: Version{Number: 3}}
	updated := &Page{ID: "123", Title: "New Title", Status: "current", Version: Version{Number: 4}}

	mux := http.NewServeMux()
	mux.HandleFunc("/rest/api/content/123", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(current)
		} else {
			// verify version was auto-incremented
			var req map[string]any
			json.NewDecoder(r.Body).Decode(&req)
			ver := req["version"].(map[string]any)
			if int(ver["number"].(float64)) != 4 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			json.NewEncoder(w).Encode(updated)
		}
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	got, err := newTestConfluenceClient(srv.URL).UpdatePage(context.Background(), "123", &UpdatePageInput{
		Title:   "New Title",
		Body:    "<p>new</p>",
		Version: 0, // auto-fetch
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Version.Number != 4 {
		t.Errorf("expected version 4, got %d", got.Version.Number)
	}
}

func TestBuildRestrictionBody(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    *RestrictionInput
		expected []map[string]any
	}{
		{
			name: "user type",
			input: &RestrictionInput{
				Operation: "update",
				Type:      "user",
				Name:      "acc-id-123",
			},
			expected: []map[string]any{
				{
					"operation": "update",
					"restrictions": map[string]any{
						"user": map[string]any{
							"type":      "known",
							"accountId": "acc-id-123",
						},
					},
				},
			},
		},
		{
			name: "group type",
			input: &RestrictionInput{
				Operation: "read",
				Type:      "group",
				Name:      "engineering",
			},
			expected: []map[string]any{
				{
					"operation": "read",
					"restrictions": map[string]any{
						"group": map[string]any{
							"type": "group",
							"name": "engineering",
						},
					},
				},
			},
		},
		{
			name: "non-user type falls into group branch",
			input: &RestrictionInput{
				Operation: "update",
				Type:      "role",
				Name:      "admins",
			},
			expected: []map[string]any{
				{
					"operation": "update",
					"restrictions": map[string]any{
						"group": map[string]any{
							"type": "group",
							"name": "admins",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := buildRestrictionBody(tt.input)
			if diff := cmp.Diff(tt.expected, got); diff != "" {
				t.Errorf("buildRestrictionBody() diff (-want +got):\n%s", diff)
			}
		})
	}
}
