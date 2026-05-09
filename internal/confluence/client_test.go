package confluence

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
