#!/usr/bin/env bash
set -Eeuo pipefail

trap cleanup SIGINT SIGTERM ERR EXIT

cleanup() {
  trap - SIGINT SIGTERM ERR EXIT
  # script cleanup here
  # Cleaning up server task
  kill $(jobs -p)
}


HOST="localhost"
PORT="4567"
HEALTH_URL="http://${HOST}:${PORT}/healthz"

./kvs -listen-port "${PORT}" &

echo "Wait for server to become ready..."
timeout --foreground -s TERM 5 bash -c ' \
  while [[ "$(curl -s -o /dev/null -L -w ''%{http_code}'' ${0})" != "204" ]]; \
  do \
    echo "Waiting for ${0}" && sleep 1; \
  done' "${HEALTH_URL}"

./kvs-cli -listen-addr "${HOST}:${PORT}"
