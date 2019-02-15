#!/usr/bin/env sh

dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

cd $dir/src
docker build -t changelog-generator .
cd ../..

docker run --rm -v $(pwd):/repo -w /repo -e GITHUB_AUTH=$GITHUB_TOKEN changelog-generator sh /app/generate.sh
