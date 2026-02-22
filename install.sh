#!/usr/bin/env bash

set -euo pipefail

function usage() {
	echo "Usage: bash install.sh [ -d directory ] [ -v version ] [ --checksums ]"
	exit 2
}

INSTALL_DIR="${HOME}/bin"
VERSION=""
VERIFY_CHECKSUMS="false"
NAME=k8spodsmetrics
OS=$(uname -o | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | tr '[:upper:]' '[:lower:]')
OS="${OS##*/}"

case "${ARCH}" in
x86_64) ARCH="amd64" ;;
aarch64 | arm64 | arm) ARCH="arm64" ;;
armv*) ARCH="arm" ;;
esac

if which go &>/dev/null; then
	gobin=$(go env GOBIN)
	gopath="$(go env GOPATH)"
	gopathbin="${gopath}/bin"
	if [ -n "${gobin}" ] && [ -d "${gobin}" ]; then
		INSTALL_DIR="${gobin}"
	elif [ -n "${gopath}" ] && [ -d "${gopathbin}" ]; then
		INSTALL_DIR="${gopathbin}"
	fi
fi

while [ $# -gt 0 ]; do
	case "$1" in
	-d)
		if [ $# -lt 2 ]; then
			usage
		fi
		INSTALL_DIR="$2"
		shift 2
		;;
	-v)
		if [ $# -lt 2 ]; then
			usage
		fi
		VERSION="$2"
		shift 2
		;;
	--checksums)
		VERIFY_CHECKSUMS="true"
		shift
		;;
	-h | --help)
		usage
		;;
	*)
		usage
		;;
	esac
done

if [ ! -d "${INSTALL_DIR}" ]; then
	echo "Directory ${INSTALL_DIR} does not exist"
	exit 1
fi
if [ ! -w "${INSTALL_DIR}" ]; then
	echo "Directory ${INSTALL_DIR} is not writable"
	exit 1
fi
APP_PATH="${INSTALL_DIR}/${NAME}"
CHECKSUMS_FILE=""

function cleanup() {
	if [ -n "${CHECKSUMS_FILE}" ] && [ -f "${CHECKSUMS_FILE}" ]; then
		rm -f "${CHECKSUMS_FILE}"
	fi
}

function verify_checksum() {
	local app_path="$1"
	local checksums_url="$2"
	local artifact_name="$3"
	local checksum_bin=""
	local expected_checksum=""
	local actual_checksum=""

	if command -v sha256sum &>/dev/null; then
		checksum_bin="sha256sum"
	elif command -v shasum &>/dev/null; then
		checksum_bin="shasum -a 256"
	else
		echo "Warning: no SHA-256 checksum tool found (sha256sum/shasum); skipping checksum verification"
		return 0
	fi

	CHECKSUMS_FILE=$(mktemp)
	echo "Downloading ${checksums_url}..."
	if ! curl -sSL --fail-with-body "${checksums_url}" -o "${CHECKSUMS_FILE}"; then
		echo "Warning: failed to download ${checksums_url}; skipping checksum verification"
		rm -f "${CHECKSUMS_FILE}"
		CHECKSUMS_FILE=""
		return 0
	fi

	expected_checksum=$(awk -v target="${artifact_name}" '$2 == target { print $1 }' "${CHECKSUMS_FILE}")
	if [ -z "${expected_checksum}" ]; then
		echo "Warning: checksum entry for ${artifact_name} not found; skipping checksum verification"
		return 0
	fi

	actual_checksum=$(${checksum_bin} "${app_path}" | awk '{ print $1 }')
	if [ "${expected_checksum}" != "${actual_checksum}" ]; then
		echo "Checksum verification failed for ${app_path}"
		echo "Expected: ${expected_checksum}"
		echo "Actual:   ${actual_checksum}"
		return 1
	fi

	echo "Checksum verified for ${app_path}"
}

# Cleanup on failure
trap 'rm -f "${APP_PATH}"' ERR
trap cleanup EXIT

echo "Installing into ${APP_PATH}..."

if [ -z "${VERSION}" ]; then
	VERSION=$(
		curl -sSL --fail-with-body https://api.github.com/repos/trezorg/${NAME}/releases/latest |
			awk -F '"' '/tag_name/ { print $4 }'
	)
	if [ -z "${VERSION}" ]; then
		echo "Failed to detect latest version"
		exit 1
	fi
fi
DOWNLOAD_URL="https://github.com/trezorg/${NAME}/releases/download/${VERSION}/${NAME}-${OS}-${ARCH}"
CHECKSUMS_URL="https://github.com/trezorg/${NAME}/releases/download/${VERSION}/checksums.txt"
echo "Downloading ${DOWNLOAD_URL}..."

if ! curl -sSL --fail-with-body "${DOWNLOAD_URL}" -o "${APP_PATH}"; then
	err=$?
	echo "Failed to download ${DOWNLOAD_URL} into ${APP_PATH}"
	exit ${err}
fi

if [ "${VERIFY_CHECKSUMS}" = "true" ]; then
	if ! verify_checksum "${APP_PATH}" "${CHECKSUMS_URL}" "${NAME}-${OS}-${ARCH}"; then
		exit 1
	fi
fi

chmod +x "${APP_PATH}"
"${APP_PATH}" --help
