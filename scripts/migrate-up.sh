#!/usr/bin/env sh
set -eu

# Uses golang-migrate/migrate
# Notes:
# - Requires `migrate` binary installed and available on PATH.
# - Applies migrations in order from migrations/.

APP_ENV=${APP_ENV:-local}

DB_URL=${DB_URL:-"postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"}

migrate -database "$DB_URL" -path ./migrations up

