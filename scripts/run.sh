#!/usr/bin/env bash

HOST="http://localhost:80/healthz"

./kvs &

echo "Wait for server to become ready..."
timeout --foreground -s TERM 5 bash -c \
  'while [[ "$(curl -s -o /dev/null -L -w ''%{http_code}'' ${0})" != "204" ]]; \
  do \
    echo "Waiting for ${0}" && sleep 1; \
  done' "${HOST}"

./kvs-cli

# Cleaning up server task
kill $(jobs -p)

