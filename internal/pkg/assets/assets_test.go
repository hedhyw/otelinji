package assets_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hedhyw/otelinj/internal/pkg/assets"
)

func TestOtelTmpl(t *testing.T) {
	t.Parallel()

	content, err := os.ReadFile("./otel.tmpl")
	require.NoError(t, err)

	assert.Equal(t, string(content), assets.OtelTmpl())
}
