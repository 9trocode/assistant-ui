# OpenClaw UI Example

This example runs a focused chat UI for an OpenClaw-style autonomous profile.

It defines an SDK-style flow named `openclaw-ui-example` and can optionally start an
embedded DevUI API server with that flow registered.

## Run

From this folder:

```bash
go run . --addr=127.0.0.1:8091 --api-base=http://127.0.0.1:7070
```

Run as a self-contained SDK example (starts embedded API + flow):

```bash
go run . --start-api --api-addr=127.0.0.1:7070 --addr=127.0.0.1:8091
```

Optional API key:

```bash
go run . --addr=127.0.0.1:8091 --api-base=http://127.0.0.1:7070 --api-key="<DEVUI_API_KEY>"
```

Then open `http://127.0.0.1:8091`.

## Docker

Build from the **repository root** so the Docker context includes `go.mod` and full source:

```bash
docker build -f examples/openclaw_ui/Dockerfile -t openclaw-ui .
docker run --rm -p 8091:8091 -p 7070:7070 openclaw-ui
```

## Standalone Project Mode

If you copy `examples/openclaw_ui` as a separate project, use the standalone template:

```bash
cp go.mod.standalone go.mod
# edit replace path to your local framework checkout
go mod tidy
go run . --start-api --api-addr=127.0.0.1:7070 --addr=127.0.0.1:8091
```

Standalone Docker (context = this folder):

```bash
docker build -f Dockerfile.standalone -t openclaw-ui-standalone .
docker run --rm -p 8091:8091 -p 7070:7070 openclaw-ui-standalone
```

Important:
- Current SDK module path is `github.com/PipeOpsHQ/agent-sdk-go/framework`.
- Use imports with `/framework/...` (for example `github.com/PipeOpsHQ/agent-sdk-go/framework/devui`).
- Importing `github.com/PipeOpsHQ/agent-sdk-go/devui` without `/framework` is not compatible with current module layout.

## Standalone Project Mode (outside monorepo)

If you copy this folder as an independent project, use the standalone module template
and standalone Dockerfile.

1) Local standalone module setup:

```bash
cp go.mod.standalone go.mod
go mod tidy
go run . --start-api --api-addr=127.0.0.1:7070 --addr=127.0.0.1:8091
```

2) Standalone Docker build (context = this folder):

```bash
docker build -f Dockerfile.standalone -t openclaw-ui-standalone .
docker run --rm -p 8091:8091 -p 7070:7070 openclaw-ui-standalone
```

Notes:
- `go.mod.standalone` uses a temporary `replace` workaround for the framework module path.
- If you want a different SDK revision, update the version in `go.mod.standalone` or set
  `--build-arg FRAMEWORK_REF=<tag-or-branch>` when using `Dockerfile.standalone`.
