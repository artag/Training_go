#!/bin/bash

OSLIST="linux windows darwin"
ARCHLIST="386 amd64 arm arm64"

for os in ${OSLIST}; do
    for arch in ${ARCHLIST}; do
        if [[ "$os/$arch" =~ ^(windows/arm64|windows/arm|darwin/arm|darwin/386)$ ]]; then continue; fi

        echo Building binary for $os $arch
        mkdir -p releases/${os}/${arch}

        CGO_ENABLED=0 GOOS=$os GOARCH=$arch \
        go build -ldflags="-s -w" -tags=inmemory \
        -o releases/${os}/${arch}

    done
done
