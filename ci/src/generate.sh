#!/usr/bin/env sh

set -e 

cp /app/package.json /repo/package.json

REVISION=$(git rev-list --tags --max-count=1)
FROM_TAG=$(git describe --tags ${REVISION})
TO_TAG=$(git describe --exact-match --tags 2> /dev/null)

if [ "$FROM_TAG" = "$TO_TAG" ]; then
    FIRST_COMMIT=$(git rev-list --max-parents=0 HEAD)
    FROM_TAG=$FIRST_COMMIT
fi

lerna-changelog --from=${FROM_TAG} --to=${TO_TAG} --next-version ${TO_TAG} > /repo/release-changelog.md
rm /repo/package.json