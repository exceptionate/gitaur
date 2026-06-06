#!/bin/bash

VERSION="0.1.2"
if [ -f .gitignore ]; then
    RELEASE="dev"
else
    RELEASE="production"
fi
COMMIT=$(git rev-parse HEAD)
DATE=$(date -u +%Y-%m-%d)
PLATFORM="$(go env GOOS)-$(go env GOARCH)"

go build -ldflags "
-X github.com/exceptionate/gitaur/cmd.Version=$VERSION
-X github.com/exceptionate/gitaur/cmd.Release=$RELEASE
-X github.com/exceptionate/gitaur/cmd.Commit=$COMMIT
-X github.com/exceptionate/gitaur/cmd.BuildDate=$DATE
-X github.com/exceptionate/gitaur/cmd.Platform=$PLATFORM
"