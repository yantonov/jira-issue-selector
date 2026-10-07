# AGENTS.md

CLI for selecting JIRA issues from the terminal (bubbletea TUI).

## Done
Done = `make check` *green*. Run `make fmt` first. Code written is not done.

## Workflow
One feature at a time, one atomic commit per logical unit. Verify *green* before starting the next feature; refactor only inside the feature at hand.

## Credentials
Read user, hostname and apikey from the OS keychain or environment variables only. The API token stays out of logs, output and files.

## Constraints
- Build via `make build` (entry point `cmd/jira-issue-selector/main.go`, output `target/jira-issue-selector`). Raise any change to either with the user first.
- Add a dependency to `go.mod` only for a clear need.

## Pointers
- Local development without a real JIRA: `make demo` (`cmd/jira-demo-server/`, `bin/start-demo.sh`).
- End-user install script: `bin/download.sh`.
