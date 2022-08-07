# otelinji

![Version](https://img.shields.io/github/v/tag/hedhyw/otelinji)
![Build Status](https://github.com/hedhyw/otelinji/actions/workflows/check.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/hedhyw/otelinji)](https://goreportcard.com/report/github.com/hedhyw/otelinji)
[![Coverage Status](https://coveralls.io/repos/github/hedhyw/otelinji/badge.svg?branch=main)](https://coveralls.io/github/hedhyw/otelinji?branch=main)

OpenTelemetry auto-instrumentation tool. It generates code with added OpenTelemetry blocks.

It injects a common open-telemetry block to all exported functions:

```go
func Example(ctx context.Context) {
    ctx, span := otel.Tracer("package").Start(ctx, "Readme")
    defer span.End()
}
```

Features:
- Custom templates are supported. Example: [template](./internal/pkg/assets/otel.tmpl).
- If the function contains a named `err` result parameter,
  then `err` will be recorded in the span.
- It supports different names of `ctx` parameter.
- This can also be used with OpenTracing, but will require a custom template.
- Using creativity, this can be used for any inserts at the beginning of the function.

## Installation

### Go

```sh
go install github.com/hedhyw/otelinji/cmd/otelinji@latest
```

## Usage

### Basic usage

Inject the layer and rewrite the file (be careful, always commit all changes first).
- `otelinji -filename input_file.go > input_file.go`

  or

- `otelinji -w -filename input_file.go`

### Recursive run

```sh
# It will inject the layer to all exported functions.
# It will ignore vendor and .git folders, test and generated files.

find . -name "*.go" \
    | grep -v "vendor/\|.git/\|_test.go" \
    | xargs -n 1 -t otelinji -w -filename
```

### Help

```
otelinji --help

Usage of otelinji:
  -filename string
        golang file [required]
  -skip-generated DO NOT EDIT
        skip files with DO NOT EDIT comment (default true)
  -template string
        path to template file [optional] (default "@/otel")
  -w    write result to file [optional]
```
