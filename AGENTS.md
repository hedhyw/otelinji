# AGENTS.md

Guidance for AI coding agents working in this repository.

## What this project is

`otelinji` (module `github.com/hedhyw/otelinji`) is an OpenTelemetry
auto-instrumentation CLI tool for Go. It parses a Go source file and injects
a tracing block (span start + deferred `EndSpanWithErr`) at the beginning of
every exported function, using AST rewriting via `github.com/dave/dst` so
comments and formatting are preserved. The injected code is rendered from a
Go text template, so with a custom template it can also be used for
OpenTracing or arbitrary inserts at the start of functions.

## CLI usage

```
otelinji -filename input_file.go        # print instrumented file to stdout
otelinji -w -filename input_file.go     # rewrite the file in place
```

Flags (parsed in `internal/pkg/config/config.go`):

- `-filename string` — Go file to process (required).
- `-w` — write the result back to the file instead of stdout.
- `-template string` — path to a custom template file; default `@/otel`
  means the embedded template `internal/pkg/assets/otel.tmpl`.
- `-skip-generated` — skip files with a `DO NOT EDIT` comment (default true).
- `-version` — print the application version and quit.

## Public Go API

- `github.com/hedhyw/otelinji/pkg/otelinji` — the only public package:
  - `otelinji.EndSpanWithErr(span, err, options...)` ends a span and, if
    `err != nil`, sets the span status to error and records the error.
  - Injected code calls this helper, so instrumented projects import this
    package. Keep it backward compatible and dependency-light (it only
    depends on `go.opentelemetry.io/otel`).

## Layout

```
cmd/otelinji/            # main entry point; `version` is set via ldflags
internal/app/            # core logic: CLI app, AST traversal, injection
internal/app/template.go # template rendering + //otelinji:check-indents handling
internal/pkg/config/     # CLI flag parsing into Config
internal/pkg/assets/     # embedded default template otel.tmpl
internal/features/       # Gherkin acceptance tests (.feature + helpers)
pkg/otelinji/            # public helper package used by injected code
```

## Template rules

The default template is `internal/pkg/assets/otel.tmpl` (embedded via
`go:embed`). Template arguments: `CtxParamName`, `FuncName`, `PackageName`,
`ReceiverType`, `IsContextUsed`, `ErrResultName`; helper function
`joinWithDot`. The `//otelinji:check-indents` comment in a template marks
lines whose presence indicates the function is already instrumented, so the
injection is skipped. The body of the template's `main` function is what
gets injected; its imports are added to target files.

`internal/features/*.exp.go.txt` / `*.in.go.txt` are test fixtures
(input/expected pairs) — update them only together with behavior changes.

## Commands

Tasks are defined in `Taskfile.yaml` (https://taskfile.dev); lint tasks
require Docker.

```sh
task check         # test + lint + lint:gherkin
task test          # go test with coverage (writes coverage.out)
task lint          # golangci-lint (dockerized, version pinned in Taskfile.yaml)
task lint:gherkin  # lint internal/features/*.feature (dockerized)
task tidy          # go mod tidy
task build         # build ./bin/otelinji with version from git describe
task run -- --help # run the CLI
```

## Conventions

- Go version: see `go.mod`. Dependencies are not vendored; run `go mod tidy`
  after changing them.
- Lint config is `.golangci.yml` (golangci-lint v2 format, almost all
  linters enabled). Run `task lint` before committing.
- Acceptance behavior is specified in Gherkin files under
  `internal/features/`; keep `.feature` files and their `_test.go`
  counterparts in sync.
- Commit messages and PR titles follow Conventional Commits (enforced on
  PR titles by the `check-pr-semantic` workflow).
- Releases are tag-driven: goreleaser + Docker image builds run from
  `.github/workflows/release.yml`; `.goreleaser.yml` and `Dockerfile`
  define packaging.
