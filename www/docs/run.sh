#!/bin/sh
set -e

RELEASES_URL="https://github.com/ecsdeployer/ecsdeployer/releases"
FILE_BASENAME="ecsdeployer"

test -z "$VERSION" && VERSION="$(curl -sfL -o /dev/null -w %{url_effective} "$RELEASES_URL/latest" |
		rev |
		cut -f1 -d'/'|
		rev)"

test -z "$VERSION" && {
	echo "Unable to get ECS Deployer version." >&2
	exit 1
}

osName=$(uname -s | tr '[:upper:]' '[:lower:]')

test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"
export TAR_FILE="$TMPDIR/${FILE_BASENAME}_${osName}_$(uname -m).tar.gz"

(
	cd "$TMPDIR"
	echo "Downloading ECS Deployer $VERSION..."
	curl -sfLo "$TAR_FILE" \
		"$RELEASES_URL/download/$VERSION/${FILE_BASENAME}_${osName}_$(uname -m).tar.gz"
	curl -sfLo "checksums.txt" "$RELEASES_URL/download/$VERSION/checksums.txt"
	echo "Verifying checksums..."
	sha256sum --ignore-missing --quiet --check checksums.txt
)

tar -xf "$TAR_FILE" -C "$TMPDIR"
"${TMPDIR}/ecsdeployer" "$@"