# Use official Go image
FROM golang:1.23

# Create and switch to a working directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod  ./
RUN go mod download

# Now copy the rest of the application
COPY . .

# Build the Go binary
RUN go build -o main .

# Expose the port the app will run on
EXPOSE 8000

# Run the application
CMD ["./main"]
