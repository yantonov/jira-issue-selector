#!/usr/bin/env sh
set -o errexit -o nounset

cd "$(dirname "$0")/.."

mkdir -p target

go build -o target/jira-issue-selector cmd/jira-issue-selector/main.go 

