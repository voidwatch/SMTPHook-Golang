FROM docker.io/library/golang:1.21-alpine AS builder

WORKDIR /app

# Copy go module files and download deps
COPY go.mod ./
COPY go.sum* ./
RUN go mod download || true

# Copy source and build
COPY . .
RUN go mod tidy && go build -o main

# --- Runtime stage ---
FROM docker.io/library/alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app

COPY --from=builder /app/main .

RUN mkdir -p /app/logs && chown -R appuser:appgroup /app

USER appuser

ENV PORT=4000

CMD ["./main"]
