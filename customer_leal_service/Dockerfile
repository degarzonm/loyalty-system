# Use an official lightweight Golang image as the base image
FROM golang:1.23-alpine AS builder

# Install dependencies (e.g., Git for Go modules hosted in private repositories)
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to take advantage of Docker caching
COPY go.mod go.sum ./

# Download dependencies only if they have changed
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application statically to reduce runtime dependencies
RUN go build -o main ./cmd

# Use a minimal base image for the final container
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8081

# Run the compiled binary
CMD ["./main"]