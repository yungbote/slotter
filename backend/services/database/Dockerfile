# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder

WORKDIR /app

# Install Git, CA certificates if needed
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Build the binary
RUN CGO_ENABLED=0 go build -o server ./cmd/server/main.go

# -----------------------------------
# Minimal final image
# -----------------------------------
FROM alpine:3.17
WORKDIR /app

# copy compiled binary from builder
COPY --from=builder /app/server /app/server

# If you have migrations or other required files, copy them too:
# COPY --from=builder /app/migrations /app/migrations

# You might also want to copy a .env if you use it. Or rely on environment in docker-compose:
# COPY .env /app/.env

EXPOSE 8080
CMD ["/app/server"]

