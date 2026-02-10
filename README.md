# OpenClaw UI Example

This example runs a focused chat UI for an OpenClaw-style autonomous profile.

It defines an SDK-style flow named `openclaw-ui-example` and can optionally start an
embedded DevUI API server with that flow registered.

## Run (monorepo)

From repo root:

```bash
go run ./examples/openclaw_ui --addr=0.0.0.0:8091 --api-base=http://127.0.0.1:7070
```

Self-contained mode (starts embedded API + flow):

```bash
go run ./examples/openclaw_ui --start-api --api-addr=0.0.0.0:7070 --addr=0.0.0.0:8091
```

Optional API key:

```bash
go run ./examples/openclaw_ui --addr=0.0.0.0:8091 --api-base=http://127.0.0.1:7070 --api-key="<DEVUI_API_KEY>"
```

Then open `http://127.0.0.1:8091`.

## Docker (monorepo)

Build from repository root so Docker context includes root `go.mod`:

```bash
docker build -f examples/openclaw_ui/Dockerfile -t openclaw-ui .
docker run --rm -p 8091:8091 -p 7070:7070 openclaw-ui
```

## Standalone Project Mode

If you copy `examples/openclaw_ui` as an independent project:

1) Create standalone module file:

```bash
cp go.mod.standalone go.mod
# edit replace path to your local SDK checkout
go mod tidy
go run . --start-api --api-addr=127.0.0.1:7070 --addr=127.0.0.1:8091
```

2) Standalone Docker build (context = this folder):

```bash
docker build -f Dockerfile.standalone -t openclaw-ui-standalone .
docker run --rm -p 8091:8091 -p 7070:7070 openclaw-ui-standalone
```

Notes:
- Imports in `main.go` use root module style:
  - `github.com/PipeOpsHQ/agent-sdk-go/devui`
  - `github.com/PipeOpsHQ/agent-sdk-go/flow`
- `go.mod.standalone` currently uses local `replace` for predictable standalone builds.
