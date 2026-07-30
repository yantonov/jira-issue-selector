#!/usr/bin/env sh
set -o errexit -o nounset

cd "$(dirname "$0")"

./build.sh

cd "$(dirname "$0")/.."

INSTALL_DIR="${HOME}/bin"

mkdir -p "${INSTALL_DIR}"

cp target/jira-issue-selector "${INSTALL_DIR}"
