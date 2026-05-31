# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gobanking-v2 ./cmd/gobanking-v2

FROM gcr.io/distroless/static-debian12
WORKDIR /
COPY --from=build /app/gobanking-v2 /gobanking-v2
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/gobanking-v2"]

