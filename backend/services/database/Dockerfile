# ----------------------
# 1) BUILD STAGE
# ----------------------
FROM golang:1.22-alpine AS build
WORKDIR /app

# Copy go.mod and go.sum first, so we can cache the module downloads.
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the source code
COPY . .

# Build the service binary
RUN go build -o /database ./cmd/database/main.go

# ----------------------
# 2) FINAL STAGE
# ----------------------
FROM alpine:3.17
WORKDIR /app

# Copy compiled binary from build stage
COPY --from=build /database /usr/local/bin/database

# Copy the migrations folder to run migrations in container if desired
COPY migrations/ /app/migrations/

# Expose port 8080 (or whatever port you set in your .env/PORT)
EXPOSE 8080

# By default, run the service
ENTRYPOINT ["database"]

