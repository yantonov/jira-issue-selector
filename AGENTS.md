# AGENTS.md

## Project
CLI tool for selecting JIRA issues from the terminal. Go 1.21, Charmbracelet huh/bubbletea, go-keyring.

## Commands
- Build: `bin/build.sh`
- Run: `bin/run.sh`
- Test: `bin/test.sh`
- Format: `bin/fmt.sh`
- Full verification: `bin/build.sh && bin/test.sh`

## Hard Constraints (MUST)
- Dependencies managed via `go.mod`. Do not add new ones without a clear need.
- All credentials (user, hostname, apikey) must only be stored via OS keychain or environment variables. Never log, print, or write the API token to a file.
- Build always via `go build -o target/jira-issue-selector ./cmd/jira-issue-selector`. Do not change the entry point or binary output path without discussion.
- Entry point: `cmd/jira-issue-selector/main.go`. Business logic in `lib/`, UI in `lib/ui/`, JIRA API client in `lib/jira/`, configuration in `lib/configuration/`.

## Definition of Done
Feature is done = `bin/build.sh` green + `bin/test.sh` green.
"Code written" is not done.

## Rules
- One feature at a time. Do not start a second until the first passes verification.
- No drive-by refactoring while the main feature is unverified.
- Before PR: run `bin/fmt.sh` and `bin/build.sh && bin/test.sh`.
- Atomic commits — one logical unit of work per commit.

## Where to Look for Details
- `lib/ui/` — TUI built with bubbletea, interface models and components.
- `lib/jira/` — JIRA REST API client.
- `lib/configuration/` — keychain and environment variable handling.
- `lib/model/` — domain types.
- `cmd/jira-demo-server/` — demo JIRA server for local development, `bin/start-demo.sh`.
- `bin/download.sh` — install script for end users.
