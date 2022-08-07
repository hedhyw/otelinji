package assets

import (
	_ "embed"
)

//go:embed otel.tmpl
var otelTmpl string

// OtelTmpl returns content of otel.tmpl.
func OtelTmpl() string {
	return otelTmpl
}
