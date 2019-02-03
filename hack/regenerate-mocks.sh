#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail


echo "Installing latest mockery..."
go get github.com/vektra/mockery/.../

echo "Generating mock implementation for interfaces..."
cd $(dirname ${BASH_SOURCE})/..
go generate ./...