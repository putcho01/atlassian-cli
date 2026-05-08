# atlassian-cli

A lightweight, native Go CLI for Atlassian Jira and Confluence Cloud (and Server/Data Center). Built as a single binary (~15 MB), with no runtime dependencies, and easy to integrate with Claude Code as a plugin skill.

## Installation

### From source (requires Go 1.23+)

```bash
go build -o atlassian-cli .
```

### go install (推奨)

`$GOPATH/bin` にインストールし、どこからでも実行できるようにします:

```bash
go install .
```

> `$GOPATH/bin` が `$PATH` に含まれていない場合は、シェル設定 (`~/.zshrc` 等) に以下を追加してください:
>
> ```bash
> export PATH="$PATH:$(go env GOPATH)/bin"
> ```

### With version info

```bash
go build -ldflags "-X github.com/putcho01/atlassian-cli/cmd.Version=1.0.0" -o atlassian-cli .
```

## Quick Start

### Cloud (atlassian.net) の場合

1. [API Token を作成](https://id.atlassian.com/manage-profile/security/api-tokens)

2. 環境変数を設定:

```bash
export JIRA_URL=https://your-domain.atlassian.net
export JIRA_EMAIL=you@example.com
export JIRA_API_TOKEN=your-api-token
```

3. 認証確認:

```bash
atlassian-cli jira myself
```

Confluence も使う場合は追加で設定:

```bash
export CONFLUENCE_URL=https://your-domain.atlassian.net/wiki
export CONFLUENCE_EMAIL=you@example.com
export CONFLUENCE_API_TOKEN=your-api-token
```

### Server/Data Center の場合

```bash
export JIRA_URL=https://jira.example.com
export JIRA_PERSONAL_TOKEN=your-pat
```

> `JIRA_EMAIL` を設定しなければ自動的に Bearer (PAT) 認証になります。

## Authentication

2つの認証方式をサポートしています:

### Cloud - API Token (Basic Auth)

Atlassian Cloud (`*.atlassian.net`) で利用。メールアドレスと API Token の組み合わせで Basic 認証します。

| Variable | Description |
|----------|-------------|
| `JIRA_URL` | Jira Cloud URL (e.g. `https://your-domain.atlassian.net`) |
| `JIRA_EMAIL` | Atlassian アカウントのメールアドレス |
| `JIRA_API_TOKEN` | [API Token](https://id.atlassian.com/manage-profile/security/api-tokens) |
| `CONFLUENCE_URL` | Confluence Cloud URL (e.g. `https://your-domain.atlassian.net/wiki`) |
| `CONFLUENCE_EMAIL` | Atlassian アカウントのメールアドレス |
| `CONFLUENCE_API_TOKEN` | API Token |

### Server/Data Center - Personal Access Token (Bearer)

自前ホスティングの Jira/Confluence Server/DC で利用。PAT で Bearer 認証します。

| Variable | Description |
|----------|-------------|
| `JIRA_URL` | Jira base URL (e.g. `https://jira.example.com`) |
| `JIRA_PERSONAL_TOKEN` | Personal Access Token |
| `CONFLUENCE_URL` | Confluence base URL |
| `CONFLUENCE_PERSONAL_TOKEN` | Personal Access Token |

> `EMAIL` が設定されていれば Basic 認証 (Cloud)、未設定なら Bearer 認証 (Server/DC) が使われます。

Only the variables for the service you use are required. For example, if you only use Jira, you don't need to set the Confluence variables.

## Commands

### Jira

```bash
# Authentication
atlassian-cli jira myself                # Show authenticated user

# Issues
atlassian-cli jira issue get PROJ-123    # Get issue (includes description)
atlassian-cli jira issue search "project = PROJ"  # Search issues via JQL
atlassian-cli jira issue create --project PROJ --summary "New task"
atlassian-cli jira issue update PROJ-123 --field summary="Updated summary"
atlassian-cli jira issue delete PROJ-123
atlassian-cli jira issue subtasks PROJ-123
atlassian-cli jira issue transition PROJ-123 "In Progress"
```

### Confluence

```bash
# Pages
atlassian-cli confluence page get 12345          # Get page content

# Labels
atlassian-cli confluence label list 12345
atlassian-cli confluence label add 12345 important,reviewed
atlassian-cli confluence label remove 12345 outdated

# Page Restrictions
atlassian-cli confluence restriction list 12345
atlassian-cli confluence restriction add 12345 --operation update --type user --name <account-id>
atlassian-cli confluence restriction remove 12345 --operation update --type user --name <account-id>
```

## Output Formats

All commands support three output formats via the `--output` / `-o` flag:

```bash
# Default: human-readable table
atlassian-cli jira issue search "project = PROJ" -o table

# Machine-readable JSON
atlassian-cli jira issue search "project = PROJ" -o json

# GitHub-flavored Markdown (great for Claude Code)
atlassian-cli jira issue search "project = PROJ" -o markdown
```

## HTML to Markdown Conversion

When using `-o markdown`, HTML content (Jira descriptions, Confluence page bodies and comments) is automatically converted to clean GitHub-flavored Markdown. Confluence storage format macros are also handled:

- Code blocks (`ac:structured-macro name="code"`) -> fenced code blocks
- Admonitions (note, info, warning, tip) -> blockquotes with labels
- Page links (`ac:link`) -> emphasized text
- Table of contents macros -> removed

## MCP Server

Start as an MCP server for AI assistant integration:

```bash
atlassian-cli mcp-server
```

### Tool Group Filtering

Filter available tools using the `--tools` flag:

```bash
# Only enable Jira issue and search tools
atlassian-cli mcp-server --tools jira_issue,jira_search

# Only enable Confluence tools
atlassian-cli mcp-server --tools confluence_page,confluence_label
```

Available tool groups:
- `jira_user` - User authentication
- `jira_issue` - Get issue, subtasks
- `jira_search` - Search issues via JQL
- `jira_create` - Create issues
- `jira_update` - Update issues
- `jira_delete` - Delete issues
- `jira_transition` - Transition issues, get available transitions
- `confluence_page` - Get page content
- `confluence_label` - Label management
- `confluence_restriction` - Page restriction management

### Claude Code Integration

Add to your Claude Code MCP settings:

```json
{
  "mcpServers": {
    "atlassian": {
      "command": "atlassian-cli",
      "args": ["mcp-server"],
      "env": {
        "JIRA_URL": "https://your-domain.atlassian.net",
        "JIRA_EMAIL": "you@example.com",
        "JIRA_API_TOKEN": "your-api-token",
        "CONFLUENCE_URL": "https://your-domain.atlassian.net/wiki",
        "CONFLUENCE_EMAIL": "you@example.com",
        "CONFLUENCE_API_TOKEN": "your-api-token"
      }
    }
  }
}
```

## Why not the official ACLI?

Atlassian provides an official CLI ([ACLI](https://developer.atlassian.com/cloud/acli/guides/introduction/)) and a [remote MCP server](https://developer.atlassian.com/cloud/acli/guides/introduction/) for AI integration. Here's how this tool differs:

### vs. ACLI

| | atlassian-cli | ACLI |
|---|---|---|
| **Distribution** | Single binary, no dependencies | Requires package installation |
| **Authentication** | Environment variables only | Interactive `acli auth login` |
| **CI/CD friendliness** | High — no browser flow needed | Limited — headless login is cumbersome |
| **Server/DC support** | Yes (PAT auth) | Cloud-focused |
| **Target audience** | Developers, automation, AI agents | Admins, bulk operations |

### vs. Atlassian remote MCP server

Atlassian's remote MCP server went GA in February 2026, but comes with trade-offs:

- **Context-heavy** — loads 73 tool schemas upfront, consuming 40–50% of the context window before any real work
- **OAuth required** — browser-based auth flow; not suitable for headless or CI/CD environments
- **Network dependency** — requires an outbound connection to Atlassian's remote server

This tool's `--tools` flag lets you load only the groups you need, keeping token usage minimal for AI agents.

### When to choose this tool

- Running in **CI/CD pipelines** or automation scripts (auth via environment variables only)
- Using both **Server/Data Center and Cloud** with a single interface
- Embedding as an **MCP server** in AI agents where context efficiency matters

## License

MIT
