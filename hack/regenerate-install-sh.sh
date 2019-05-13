#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

echo "Installing latest godownloader..."
go get github.com/goreleaser/godownloader

INSTALL_SCRIPT="install.sh"

echo "Generating install.sh script..."
cd $(dirname ${BASH_SOURCE})/..
godownloader .goreleaser.yml --repo=pkosiec/terminer > ${INSTALL_SCRIPT}
chmod +x ${INSTALL_SCRIPT}