#!/bin/bash

[[ -n "$GOOS" && -n "$GOARCH" ]] && OUT_DIR="bin/$GOOS/$GOARCH" || OUT_DIR="bin/"

[[ -d "$OUT_DIR" ]] || mkdir -p "$OUT_DIR"

CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o "$OUT_DIR" ./...
