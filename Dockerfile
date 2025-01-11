# Use the official Golang image as a builder
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the rest of the application
COPY backend/ ./

# Build the Go application
RUN go build -o scraper ./internal/scraper

# Use a minimal image for the final container
FROM debian:bullseye-slim

# Set the working directory in the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/scraper .

# Expose a port (if necessary for your application)
EXPOSE 8080

# Run the application
CMD ["./scraper"]
