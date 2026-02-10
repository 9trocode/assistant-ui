# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY . .

ARG FRAMEWORK_REPO=https://github.com/PipeOpsHQ/agent-sdk-go.git
ARG FRAMEWORK_REF=main

RUN set -eux; \
    if [ -f go.mod ] && [ -f examples/openclaw_ui/main.go ]; then \
      go mod download; \
      CGO_ENABLED=0 GOOS=linux go build -o /app/openclaw-ui ./examples/openclaw_ui; \
    elif [ -f main.go ]; then \
      rm -rf /tmp/framework-src; \
      git clone --depth 1 --branch "${FRAMEWORK_REF}" "${FRAMEWORK_REPO}" /tmp/framework-src; \
      MOD_DIR="/tmp/framework-src"; \
      if [ -f /tmp/framework-src/framework/go.mod ]; then MOD_DIR="/tmp/framework-src/framework"; fi; \
      go mod init openclaw-ui-local || true; \
      go mod edit -replace github.com/PipeOpsHQ/agent-sdk-go/framework=${MOD_DIR}; \
      go mod tidy; \
      CGO_ENABLED=0 GOOS=linux go build -o /app/openclaw-ui .; \
    else \
      echo "Unsupported build context. Use repo root or examples/openclaw_ui."; \
      exit 1; \
    fi

# Runtime stage
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/openclaw-ui /app/openclaw-ui

# OpenClaw chat UI + embedded DevUI API
EXPOSE 8091 7070

ENTRYPOINT ["/app/openclaw-ui"]
CMD ["--start-api", "--api-addr=127.0.0.1:7070", "--addr=127.0.0.1:8091"]
