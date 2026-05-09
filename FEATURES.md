# atlassian-cli Feature List & Comparison

## Our atlassian-cli

A lightweight, native Go CLI for Atlassian Jira and Confluence Cloud (and Server/Data Center). Built as a single binary, and easy to integrate with Claude Code as a plugin skill.

## Jira Features

| Feature | Command | Status |
|---------|---------|--------|
| Auth verification | `jira myself` | Done |
| Get issue | `jira issue get <key>` | Done |
| Open issue in browser | `jira issue open <key>` | Done |
| Search issues (JQL) | `jira issue search <jql>` | Done |
| Create issue | `jira issue create --project --summary ...` | Done |
| Update issue | `jira issue update <key> --field value` | Done |
| Delete issue | `jira issue delete <key>` | Done |
| Subtasks | `jira issue subtasks <key>` | Done |
| Transition issue | `jira issue transition <key> <status>` | Done |
| List comments | `jira issue comment list <key>` | Done |
| Add comment | `jira issue comment add <key> <body>` | Done |

## Confluence Features

| Feature | Command | Status |
|---------|---------|--------|
| Get page | `confluence page get <id>` | Done |
| List labels | `confluence label list <id>` | Done |
| Add labels | `confluence label add <id> <labels>` | Done |
| Remove label | `confluence label remove <id> <label>` | Done |
| List restrictions | `confluence restriction list <id>` | Done |
| Add restriction | `confluence restriction add <id> --name ...` | Done |
| Remove restriction | `confluence restriction remove <id> --name ...` | Done |

## General Features

| Feature | Details | Status |
|---------|---------|--------|
| Cloud API Token auth (Basic) | Via environment variables | Done |
| Bearer token auth (PAT) | Via environment variables | Done |
| Output formats (table/json/markdown) | `--output` / `-o` flag | Done |
| HTML -> Markdown conversion | Confluence macros supported | Done |
| MCP server (`mcp-server`) | stdio JSON-RPC | Done |
| MCP tool group filtering (`--tools`) | 10 groups | Done |
| Version command | `version` | Done |

## Architecture

| Aspect | atlassian-cli |
|--------|---------------|
| Language | Go |
| Distribution | Single binary |
| Runtime dependency | None |
| Integration model | CLI + MCP server + Claude Code skill |
| License | MIT |
