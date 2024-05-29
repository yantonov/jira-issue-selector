#!/usr/bin/env sh
set -o errexit -o nounset

cd "$(dirname "$0")/.."

go build -o target/jira-issue-selector cmd/jira-issue-selector/main.go 

