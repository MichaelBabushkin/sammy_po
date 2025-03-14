# 1. Builder
FROM golang:1.21beta1 as builder
WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN go build -o scraper ./internal/scraper

# 2. Minimal runtime
FROM debian:bullseye-slim
WORKDIR /app

COPY --from=builder /app/scraper .
EXPOSE 8080
CMD ["./scraper"]
