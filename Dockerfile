FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum* ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o findarr ./cmd/server

# Use a minimal alpine image for the final container
FROM alpine:latest

WORKDIR /app

# Install ca-certificates and wget for HTTPS requests and healthchecks
RUN apk --no-cache add ca-certificates wget

# Copy the binary from the builder stage
COPY --from=builder /app/findarr .

# Copy web templates and static files
COPY --from=builder /app/web ./web

# Expose the port the server listens on
EXPOSE 8080

# Run the application
CMD ["./findarr"]

