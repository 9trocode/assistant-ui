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
