FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-server ./cmd/api/main.go

# Use a lightweight alpine image for the runtime
FROM alpine:latest

WORKDIR /
COPY --from=builder /api-server /api-server
COPY .env.example .env

EXPOSE 8080

CMD ["/api-server"]
