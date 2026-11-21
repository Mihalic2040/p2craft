# Build stage
FROM golang:1.27.2-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Compile
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /hole \
    ./cmd/hole

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /hole /usr/local/bin/hole

ENTRYPOINT ["hole"]
