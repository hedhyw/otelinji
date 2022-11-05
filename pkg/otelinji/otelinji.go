package otelinji

import (
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// EndSpanWithErr finishes the span and records error.
func EndSpanWithErr(
	span trace.Span,
	err error,
	options ...trace.SpanEndOption,
) {
	if err != nil {
		span.SetStatus(codes.Error, "error")
		span.RecordError(err)
	}

	span.End(options...)
}
