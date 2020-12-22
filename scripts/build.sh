#!/usr/bin/env bash

go get ./...

go build ./cmd/kvs
go build ./cmd/kvs-cli


