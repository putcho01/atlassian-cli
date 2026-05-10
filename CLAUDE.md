# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development

```bash
# Build
go build -o atlassian-cli .

# Run tests
go test --shuffle on --race -v ./...

# Run tests for a specific package
go test ./internal/htmlconv/...

# Run a single test
go test ./internal/htmlconv/ -run TestConvert

# Lint
golangci-lint run
```

## Architecture

The project is a single Go binary (`main.go` → `cmd.Execute()`) structured as a CLI + MCP server.

### Package layout

- **`cmd/`** — Cobra command tree. Each file corresponds to a command group (e.g., `jira_issue_get.go`, `confluence_page.go`). Commands instantiate an API client, call it, and delegate rendering to `internal/formatter`.
- **`internal/config/`** — Reads `JIRA_*` / `CONFLUENCE_*` environment variables. If `EMAIL` is set, Basic auth (Cloud) is used; otherwise Bearer auth (Server/DC). `LoadJiraConfig()` / `LoadConfluenceConfig()` are the entry points.
- **`internal/httpclient/`** — Thin HTTP client wrapper that injects auth headers.
- **`internal/jira/`** — Jira REST API client and type definitions (`types.go`).
- **`internal/confluence/`** — Confluence REST API client and type definitions (`types.go`).
- **`internal/formatter/`** — Renders API responses in three formats: `table` (tab-aligned), `json`, `markdown`. HTML content is piped through `htmlconv.Convert()` when using markdown format.
- **`internal/htmlconv/`** — HTML-to-GFM converter. Handles standard HTML tags plus Confluence storage format macros (`ac:structured-macro`, `ac:link`, `ac:emoticon`, etc.). This is where macro rendering logic lives.
- **`internal/mcptools/`** — MCP server implementation using `github.com/modelcontextprotocol/go-sdk`. `groups.go` maps group names to tool names; `server.go` creates the server and registers tools via `maybeAdd()`.
- **`internal/tui/`** — Bubbletea-based interactive issue picker (activated with `jira issue search -i`).

### Key design patterns

- **Output format**: All commands take `--output` / `-o` (table/json/markdown) via a persistent root flag. Commands read it with `cmd.Root().PersistentFlags().GetString("output")`, then call `formatter.ParseFormat()`.
- **Auth dispatch**: `config.loadConfig()` tries `_API_TOKEN` first, then `_PERSONAL_TOKEN`. If `Email` is non-empty in the resulting struct, the HTTP client uses Basic auth; otherwise Bearer.
- **MCP tool filtering**: `mcptools.maybeAdd()` checks `isToolEnabled(name, enabledGroups)` before registering each tool. The `--tools` CLI flag passes a comma-separated list of group names.
- **HTML conversion**: `htmlconv.Convert()` uses `golang.org/x/net/html` to parse, then walks the node tree with `renderElement()`. Confluence-specific tags (`ac:*`, `ri:*`) are handled before standard HTML tags.

## Release

Releases are automated via GitHub Actions (`.github/workflows/release.yml`) triggered on `v*` tags. GoReleaser builds cross-platform binaries for linux/darwin/windows × amd64/arm64.
