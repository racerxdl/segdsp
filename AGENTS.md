# SegDSP

Go SDR demodulator. Single-module, single-binary project.

## Commands

```bash
go build -o segdsp .
go test -v -race ./...
golangci-lint run
go vet ./...
go fmt ./...
```

No Makefile or task runner exists.

## Architecture

- Entry point: `segdsp.go` (package main, all root `*.go` files)
- `dsp/` — DSP primitives (FIR, FFT, NCO, AGC, resamplers). Contains arch-specific code:
  - `dsp/native/amd64/` — AMD64 assembly optimizations
  - `dsp/native/c/` — C fallback implementations
- `demodcore/` — Demodulator interface + FM/AM implementations
- `eventmanager/` — Channel-based pub/sub event bus
- `recorders/` — WAV file recording
- `content/` — Static web UI served by the binary

## Gotchas

- Docker builds use `CGO_ENABLED=0` for static binaries despite native C code existing
- arm32v6 is cross-compiled on the host (not in Docker) due to a Go compiler bug in Alpine — see `multi-build.sh`
- Binary is named `segdsp` locally but `segdsp_worker` inside Docker containers
- The app **cannot start** without a running `radioserver` instance (external SDR IQ source, not in this repo)
- Default branch is `master`; releases are tag-triggered via Travis CI

## Testing

All tests are pure unit tests (DSP math primitives). No integration tests, no external service dependencies, no test fixtures. Benchmark files exist (`benchmark_*.go`).

## Conventions

- HTTP routing uses stdlib `http.HandleFunc` directly, no router library
- CLI flags parsed with stdlib `flag`, env vars set in `common.go:setEnv()`
- WebSocket endpoint at `/ws`, static UI at `/` and `/static/*`
