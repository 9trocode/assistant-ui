# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/openclaw-ui .

# Runtime stages
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/openclaw-ui /app/openclaw-ui

# Build-time args -> runtime env defaults
ARG OPENCLAW_START_API=true
ARG OPENCLAW_API_ADDR=0.0.0.0:7070
ARG OPENCLAW_ADDR=0.0.0.0:8091
ARG OPENCLAW_API_BASE=
ARG OPENCLAW_API_KEY=

# Runtime configuration
# - OPENCLAW_START_API: "true" to start embedded DevUI API
# - OPENCLAW_API_ADDR: embedded API listen address
# - OPENCLAW_ADDR: OpenClaw UI listen address
# - OPENCLAW_API_BASE: optional external DevUI API URL (used when not embedding)
# - OPENCLAW_API_KEY: optional DevUI API key sent as X-API-Key
ENV OPENCLAW_START_API=${OPENCLAW_START_API} \
    OPENCLAW_API_ADDR=${OPENCLAW_API_ADDR} \
    OPENCLAW_ADDR=${OPENCLAW_ADDR} \
    OPENCLAW_API_BASE=${OPENCLAW_API_BASE} \
    OPENCLAW_API_KEY=${OPENCLAW_API_KEY}

# OpenClaw chat UI + embedded DevUI API
EXPOSE 8091 7070

ENTRYPOINT ["/bin/sh", "-c"]
CMD ["set -eu; args=\"--addr=${OPENCLAW_ADDR}\"; if [ \"${OPENCLAW_START_API}\" = \"true\" ] || [ \"${OPENCLAW_START_API}\" = \"1\" ]; then args=\"$args --start-api --api-addr=${OPENCLAW_API_ADDR}\"; fi; if [ -n \"${OPENCLAW_API_BASE:-}\" ]; then args=\"$args --api-base=${OPENCLAW_API_BASE}\"; fi; if [ -n \"${OPENCLAW_API_KEY:-}\" ]; then args=\"$args --api-key=${OPENCLAW_API_KEY}\"; fi; exec /app/openclaw-ui $args"]
