#!/usr/bin/env sh
set -o errexit -o nounset

cd "$(dirname "$0")/.."

go build -o target/jira-demo-server cmd/jira-demo-server/main.go