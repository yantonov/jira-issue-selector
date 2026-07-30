#!/usr/bin/env sh
set -o errexit -o nounset

cd "$(dirname "$0")/.."

INSTALL_DIR="${HOME}/bin"
mkdir -p "${INSTALL_DIR}"

BUILD_ARTIFACT="target/jira-issue-selector"

if [ -f "${BUILD_ARTIFACT}" ]; then
    cp target/jira-issue-selector "${INSTALL_DIR}"
else
    echo "binary artifact not found"
    echo "build it first: ./bin/build.sh"
    exit 1
fi
