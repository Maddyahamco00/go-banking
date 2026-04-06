# Go Banking System (In Progress)

A backend banking system built with Go, PostgreSQL, and Docker.

## Features (Planned)
- Account creation
- Deposit & withdrawal
- Money transfer
- Transaction logging

## Tech Stack
- Go (Golang)
- PostgreSQL
- Docker
- golang-migrate

## Current Progress
- Database containerized using Docker
- Migration system implemented
- Accounts table created
- API layer in progress

## Run Locally

Start PostgreSQL:
docker run --name postgres-banking -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=banking -p 5433:5432 -d postgres

Run migrations:
migrate -path migrations -database "postgres://admin:secret@localhost:5433/banking?sslmode=disable" up