package mcptools

var toolGroups = map[string][]string{
	"jira_user":              {"jira_get_myself"},
	"jira_issue":             {"jira_get_issue", "jira_get_subtasks"},
	"jira_search":            {"jira_search_issues"},
	"jira_create":            {"jira_create_issue"},
	"jira_update":            {"jira_update_issue"},
	"jira_delete":            {"jira_delete_issue"},
	"jira_transition":        {"jira_transition_issue", "jira_get_transitions"},
	"confluence_page":        {"confluence_get_page"},
	"confluence_label":       {"confluence_list_labels", "confluence_add_labels", "confluence_remove_label"},
	"confluence_restriction": {"confluence_list_restrictions", "confluence_add_restriction", "confluence_remove_restriction"},
}

func isToolEnabled(name string, enabledGroups []string) bool {
	if len(enabledGroups) == 0 {
		return true
	}
	for _, group := range enabledGroups {
		tools, ok := toolGroups[group]
		if !ok {
			continue
		}
		for _, t := range tools {
			if t == name {
				return true
			}
		}
	}
	return false
}

// AllGroups returns all available tool group names.
func AllGroups() []string {
	groups := make([]string, 0, len(toolGroups))
	for g := range toolGroups {
		groups = append(groups, g)
	}
	return groups
}
