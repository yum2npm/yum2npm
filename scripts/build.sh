#!/bin/bash

[[ -n "$GOOS" && -n "$GOARCH" ]] && OUT_DIR="bin/$GOOS/$GOARCH" || OUT_DIR="bin/"

[[ -d "$OUT_DIR" ]] || mkdir -p "$OUT_DIR"

[[ -n "$CI_COMMIT_TAG" ]] && VERSION="$CI_COMMIT_TAG" || VERSION="devel"

CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static" -X main.Version="'"$VERSION"'"' -o "$OUT_DIR" ./...
