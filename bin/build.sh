#!/usr/bin/env sh
set -o errexit -o nounset

cd "$(dirname "$0")/.."

mkdir -p target

# the package is built, not the file list,
# otherwise the go tool does not stamp the commit into the binary
go build -o target/jira-issue-selector ./cmd/jira-issue-selector 

