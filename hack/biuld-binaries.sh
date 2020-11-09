#!/bin/bash

set -e -x -u

BUILD_VALUES= ./hack/build.sh

GOOS=darwin GOARCH=amd64 go build -o github-contrib-darwin-amd64 ./cmd/...
GOOS=linux GOARCH=amd64 go build -o github-contrib-linux-amd64 ./cmd/...
GOOS=windows GOARCH=amd64 go build -o github-contrib-windows-amd64.exe ./cmd/...

shasum -a 256 ./github-contrib-*-amd64*
