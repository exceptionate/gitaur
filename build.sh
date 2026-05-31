#!/bin/bash

VERSION="0.1.0"
RELEASE="dev"
COMMIT=$(git rev-parse HEAD)
DATE=$(date -u +%Y-%m-%d)
PLATFORM="$(go env GOOS)-$(go env GOARCH)"

go build -ldflags "
-X 'github.com/exceptionate/gitaur/version.Version=$VERSION'
-X 'github.com/exceptionate/gitaur/version.Release=$RELEASE'
-X 'github.com/exceptionate/gitaur/version.Commit=$COMMIT'
-X 'github.com/exceptionate/gitaur/version.BuildDate=$DATE'
-X 'github.com/exceptionate/gitaur/version.Platform=$PLATFORM'
"