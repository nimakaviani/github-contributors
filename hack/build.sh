#!/bin/bash

set -e -x -u

go fmt ./cmd/... ./pkg/...

# build without website assets
go build -o github-contrib ./cmd/...
./github-contrib version
