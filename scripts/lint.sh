#!/bin/bash

shopt -s extglob

test -z "$(gofmt -s -l !(.go)/)" && exit 0 || exit 1
staticcheck ./... || exit 1
