#!/usr/bin/env sh
set -eu

DB_URL=${DB_URL:-"postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"}

migrate -database "$DB_URL" -path ./migrations down 1

