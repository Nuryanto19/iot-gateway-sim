#!/bin/bash
set -e

# Use token file for best practice, ensure token is match with environtment config in docker compose.yaml
influx apply --force yes \
  -token "my-super-secret-admin-token-123" \
  -o "my-iot-org" \
  -f /docker-entrypoint-initdb.d/dashboard.yaml
