# Build React frontend
FROM node:18 AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Build Go backend
FROM golang:1.23
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/frontend/build ./frontend/build
RUN go build -o main .

EXPOSE 8000
CMD ["./main"]
