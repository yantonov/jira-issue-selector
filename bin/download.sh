#!/usr/bin/env sh

set -eu

SCRIPT="$(basename "$0")"
cd "$(dirname "$0")"

APP_NAME="jira-issue-selector"

# Detect OS
case "$(uname -s)" in
  Linux*)
    OS="linux"
    ARCHIVE_EXT="tar.gz"
    ;;
  Darwin*)
    OS="macos"
    ARCHIVE_EXT="tar.gz"
    ;;
  MINGW*|MSYS*|CYGWIN*)
    OS="windows"
    ARCHIVE_EXT="zip"
    ;;
  *)
    echo "Unsupported OS: $(uname -s)"
    exit 1
    ;;
esac

REPO="yantonov/jira-issue-selector"

# Get latest tag
LATEST_TAG=$(
  curl -fsSL "https://api.github.com/repos/${REPO}/tags" \
  | jq -r '.[0].name'
)

ARCH="amd64"

ALIAS_APP_NAME="jira-issue-selector"
ARCHIVE_NAME="${ALIAS_APP_NAME}-${LATEST_TAG}-${OS}-${ARCH}.${ARCHIVE_EXT}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${ARCHIVE_NAME}"

echo "Latest tag: ${LATEST_TAG}"
echo "Downloading: ${DOWNLOAD_URL}"

TMP_DIR="$(mktemp -d)"
ARCHIVE_PATH="${TMP_DIR}/${ALIAS_APP_NAME}.${ARCHIVE_EXT}"

echo "${ARCHIVE_PATH}"

# Download archive
curl -fL "${DOWNLOAD_URL}" -o "${ARCHIVE_PATH}"

# Extract archive
if [ "${ARCHIVE_EXT}" = "zip" ]; then
  unzip -q "${ARCHIVE_PATH}" -d "${TMP_DIR}"
else
  tar -xzf "${ARCHIVE_PATH}" -C "${TMP_DIR}"
fi

# Find binary inside extracted files
BIN_PATH="$(find "${TMP_DIR}" -type f -exec sh -c 'test -x "$1"' _ {} \; -print | head -n 1)"

if [ -z "${BIN_PATH}" ]; then
  echo "Executable ${ALIAS_APP_NAME} is not found in the archive ${TMP_DIR}"
  rm -rf "${TMP_DIR}"
  exit 1
fi

TARGET_DIR="${HOME}/bin"
mkdir -p "${TARGET_DIR}"

# Copy binary to the target directory
cp "${BIN_PATH}" "${TARGET_DIR}/${APP_NAME}"
chmod +x "${TARGET_DIR}/${APP_NAME}"

# Cleanup
rm -rf "${TMP_DIR}"

echo "Installed: ${TARGET_DIR}/${APP_NAME}"
