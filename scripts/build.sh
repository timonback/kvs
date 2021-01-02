#!/usr/bin/env bash

go get -v -t -d ./...

go build ./cmd/kvs
go build ./cmd/kvs-cli


