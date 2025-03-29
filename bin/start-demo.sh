#!/usr/bin/env sh
set -o errexit -o nounset

cd "$(dirname "$0")/../target"

JIRA_USER=user JIRA_API_KEY=key JIRA_HOSTNAME=http://localhost:8090 ./jira-issue-selector


