package otelinji_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/hedhyw/otelinj/pkg/otelinji"
)

func TestEndSpanWithErrOK(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var buf bytes.Buffer

	stdExporter, err := stdouttrace.New(stdouttrace.WithWriter(&buf))
	require.NoError(t, err)

	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(stdExporter))

	_, span := tracerProvider.Tracer("otelinji_test").Start(ctx, "TestEndSpanWithErrOK")
	otelinji.EndSpanWithErr(span, nil)

	err = tracerProvider.ForceFlush(ctx)
	require.NoError(t, err)

	err = tracerProvider.Shutdown(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, buf.String())
}

func TestEndSpanWithErrFailure(t *testing.T) {
	t.Parallel()

	const errTest semerr.Error = "example error"

	ctx := context.Background()

	var buf bytes.Buffer

	stdExporter, err := stdouttrace.New(stdouttrace.WithWriter(&buf))
	require.NoError(t, err)

	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(stdExporter))

	_, span := tracerProvider.Tracer("otelinji_test").Start(ctx, "TestEndSpanWithErrFailure")
	otelinji.EndSpanWithErr(span, errTest)

	err = tracerProvider.ForceFlush(ctx)
	require.NoError(t, err)

	err = tracerProvider.Shutdown(ctx)
	require.NoError(t, err)
	require.Contains(t, buf.String(), errTest.Error())
}
