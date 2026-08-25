#!/usr/bin/env sh
set -o errexit -o nounset

cd "$(dirname "$0")/.."

INSTALL_DIR="${HOME}/.local/bin"

mkdir -p "${INSTALL_DIR}"

BUILD_ARTIFACT="target/jira-issue-selector"

./bin/build.sh

cp "${BUILD_ARTIFACT}" "${INSTALL_DIR}"
echo "Installed to ${INSTALL_DIR}"
