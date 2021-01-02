#!/usr/bin/env bash
set -Eeuo pipefail

go get -v -t -d ./...

go build ./cmd/kvs
go build ./cmd/kvs-cli
