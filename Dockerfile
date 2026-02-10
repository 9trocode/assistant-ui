# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/openclaw-ui ./examples/openclaw_ui

# Runtime stage
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/openclaw-ui /app/openclaw-ui

# OpenClaw chat UI + embedded DevUI API
EXPOSE 8091 7070

ENTRYPOINT ["/app/openclaw-ui"]
CMD ["--start-api", "--api-addr=127.0.0.1:7070", "--addr=127.0.0.1:8091"]
