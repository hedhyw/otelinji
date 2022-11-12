package features

import (
	"bytes"
	"go/format"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hedhyw/otelinji/internal/app"
	"github.com/hedhyw/otelinji/internal/pkg/config"
)

// RequireRunApp is a test helper that runs the application with provided
// CLI arguments.
func RequireRunApp(tb testing.TB, args ...string) string {
	tb.Helper()

	cfg, err := config.FromCLI(args, "v0.0.1")
	require.NoError(tb, err)

	var buf bytes.Buffer

	err = app.New(cfg).Run(&buf)
	require.NoError(tb, err)

	return buf.String()
}

// RequireFormatGoSource formats go source code using `format.Source`.
func RequireFormatGoSource(tb testing.TB, input string) string {
	tb.Helper()

	data, err := format.Source([]byte(input))
	require.NoError(tb, err)

	return string(data)
}

// RequireTempFile creates a temporary file with a content. The file
// will be remove at cleanup step.
func RequireTempFile(tb testing.TB, content []byte) (path string) {
	tb.Helper()

	f, err := os.CreateTemp(tb.TempDir(), "")
	require.NoError(tb, err)

	path = f.Name()

	tb.Cleanup(func() { assert.NoError(tb, os.Remove(path)) })

	defer func() { require.NoError(tb, f.Close()) }()

	_, err = f.Write(content)
	require.NoError(tb, err)

	return path
}
