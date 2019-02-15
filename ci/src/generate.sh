#!/usr/bin/env sh

cp /app/package.json /repo/package.json
lerna-changelog --next-version $(git describe --exact-match --tags 2> /dev/null) > /repo/release-changelog.md
rm /repo/package.json