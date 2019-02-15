#!/usr/bin/env bash

set -e

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})

GITHUB_AUTH=$GITHUB_TOKEN
$SCRIPT_ROOT/generate-changelog.sh

CHANGELOG_PATH=$(pwd)/release-changelog.md
cat $CHANGELOG_PATH

curl -sL https://git.io/goreleaser | bash -s -- --release-notes $CHANGELOG_PATH