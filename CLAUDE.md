# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

GenAi-Wrapper is a private Go SDK maintained by the Checkmarx AI Squad that abstracts calls to LLM providers (OpenAI, Checkmarx's internal CxOne AI gateway, and a LiteLLM proxy). It provides stateless and stateful chat wrappers, conversation history persistence, and secret-masking of prompts/responses before they're sent or logged.

## Commands

```bash
go build ./...                                  # build everything
go test ./...                                   # run all tests
go test ./pkg/wrapper/...                       # run tests for a single package
go test ./pkg/wrapper/ -run TestName -v         # run a single test
go test ./... -coverpkg=./... -coverprofile cover.out   # coverage (matches CI)
go vet ./...
golangci-lint run -c .golangci.yml --timeout 10m  # matches the Lint GitHub Action (v2.11.3)
go mod tidy                                     # required before lint in CI
govulncheck ./...                               # matches the Govulncheck GitHub Action
```

CI (GitHub Actions, `.github/workflows/`) runs lint, govulncheck, a Codecov coverage scan (`go test ./... -coverpkg=./...`), and a Checkmarx One scan (`cx-one-scan.yaml`) on every PR. `pr-linter.yaml` enforces PR title/metadata conventions. There is no Makefile; use `go` and `golangci-lint` directly.

## Architecture

The codebase is organized as three layers, from low-level HTTP transport up to the public API:

1. **`internal/`** — raw HTTP transport to LLM providers, hidden from consumers.
   - `internal/genaiExternal.go` (`WrapperImpl`) — generic OpenAI-compatible chat completion caller. Handles Checkmarx-gateway auth (`cxAuth` + `MetaData` headers: `X-Request-ID`, `X-Tenant-ID`, `X-Feature`, etc.) vs. direct OpenAI auth (plain API key), and injects `setupMessages` (system/developer prompts registered via `SetupCall`) into the message list at the right position depending on model (GPT-4 vs. others, which insert after the last user message).
   - `internal/litellm_wrapper.go` (`LitellmWrapper`) — separate, simpler transport for the LiteLLM proxy; always uses Bearer auth + `MetaData` headers, no setup-message injection.
   - `internal/gpt.go` — shared request/response types (`ChatCompletionRequest/Response`, `ErrorResponse`) and the `Wrapper` interface both transports implement; also `NewWrapperFactory` / `NewLitellmWrapperFactory`.
   - `internal/codes.go` — provider error code / finish-reason constants (e.g. `context_length_exceeded`, `FinishReasonLength`) used to trigger history truncation.
   - `internal/secrets/` — regex+entropy-based secret detector (`MaskSecrets`), rules loaded from an embedded `regex_rules.json`. Used to redact secrets from conversation content before it's sent upstream or surfaced back.
   - `internal/api/redirect_prompt/` — generated protobuf/gRPC code (not yet wired into any wrapper).

2. **`pkg/wrapper/`** — public wrapper implementations built on `internal.Wrapper`:
   - `StatelessWrapper` (`stateless_wrapper.go`) — takes explicit `history []message.Message` on every call, masks secrets in the full conversation, enforces an optional user-message `limit`, and on a `FinishReasonLength` response recursively retries after dropping `dropLen` oldest messages from history.
   - `StatefulWrapper` (`stateful_wrapper.go`) — wraps a `StatelessWrapper` plus a `connector.Connector`; looks up history by `uuid.UUID`, delegates the call, then appends and persists the updated history. `NewStatefulWrapper` (non-`New` variant) is deprecated in favor of `NewStatefulWrapperNew`.
   - `LitellmWrapper` (`litellm_wrapper.go`) — thin pass-through to `internal.LitellmWrapper`, defaults to `models.DefaultModel` if none given.

3. **`pkg/connector/`** — history persistence abstraction (`Connector` interface: `HistoryById` / `SaveHistory` / `DeleteHistory`). `FileSystemConnector` is the only implementation, storing one JSON file per conversation UUID under `<baseDir>/cx-gpt/`. It defends against path traversal via `safeBasePath`/`validatePath`, which resolve and confirm paths stay inside the base directory before any read/write.

4. **`pkg/message/`, `pkg/role/`, `pkg/models/`, `pkg/maskedSecret/`** — shared value types: `Message`/`ChatResponse`/`MetaData`/`TokenUsage`, role constants (`system`/`assistant`/`user`/`developer`), model name constants (OpenAI + Claude), and the masked-secret result type.

5. **`example/`** — standalone `main` demonstrating the SDK as a CLI chat tool (`example/main.go` + `cxoneai.go`/`openai.go`/`utils.go` for provider-specific key/endpoint lookup). Currently only wires up the LiteLLM path end-to-end; OpenAI/CxOne key-fetching helpers exist but `CallAIandPrintResponse` only supports `-ai LiteLLM`.

### Key request flow

`StatefulWrapper.SecureCall` → loads history via `Connector` → `StatelessWrapper.SecureCallReturningFullResponse` (masks secrets, builds `ChatCompletionRequest`) → `internal.Wrapper.Call` (injects setup messages, sets auth headers, does the HTTP POST) → on `context_length_exceeded` / `FinishReasonLength`, the call recurses with `dropLen` messages trimmed from the front of history → response is appended to history and persisted.

### Auth model

Two auth modes distinguished by whether `MetaData` is nil:
- `MetaData != nil` → Checkmarx gateway mode: `Authorization: Bearer <cxAuth>` plus `X-Request-ID`/`X-Tenant-ID`/`User-Agent`/`X-Feature` headers, and gateway-specific error handling via the `X-Gen-Ai-ErrorCode` response header.
- `MetaData == nil` → direct OpenAI mode: `Authorization: Bearer <apiKey>` (the wrapper's own configured key), standard OpenAI error body parsing.
