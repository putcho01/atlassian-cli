package config

import (
	"fmt"
	"os"
)

// AtlassianConfig holds connection settings for an Atlassian service.
// Email が設定されていれば Basic 認証 (Cloud)、未設定なら Bearer 認証 (Server/DC)。
type AtlassianConfig struct {
	URL            string
	Email          string
	Token          string
	DefaultProject string
}

func loadConfig(prefix string) (*AtlassianConfig, error) {
	url := os.Getenv(prefix + "_URL")
	email := os.Getenv(prefix + "_EMAIL")
	token := os.Getenv(prefix + "_API_TOKEN")
	if token == "" {
		token = os.Getenv(prefix + "_PERSONAL_TOKEN")
	}

	if url == "" {
		return nil, fmt.Errorf("%s_URL environment variable is required", prefix)
	}
	if token == "" {
		return nil, fmt.Errorf("%s_API_TOKEN (or %s_PERSONAL_TOKEN) environment variable is required", prefix, prefix)
	}

	return &AtlassianConfig{
		URL:            url,
		Email:          email,
		Token:          token,
		DefaultProject: os.Getenv(prefix + "_DEFAULT_PROJECT"),
	}, nil
}

func LoadJiraConfig() (*AtlassianConfig, error)       { return loadConfig("JIRA") }
func LoadConfluenceConfig() (*AtlassianConfig, error)  { return loadConfig("CONFLUENCE") }
