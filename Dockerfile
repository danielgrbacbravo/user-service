# Start from the official Golang image for building
FROM golang:1.24 as builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/api/main.go

# Use a minimal image for running
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder
COPY --from=builder /app/user-service .

EXPOSE 8080

CMD ["./user-service"]
