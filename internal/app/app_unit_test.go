package app

import (
	"testing"

	"github.com/hedhyw/otelinji/internal/pkg/assets"

	"github.com/stretchr/testify/assert"
)

func TestGetTemplateContent(t *testing.T) {
	t.Parallel()

	t.Run("internal", func(t *testing.T) {
		t.Parallel()

		content, err := getTemplateContent("@/otel")
		if assert.NoError(t, err) {
			assert.Equal(t, assets.OtelTmpl(), content)
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		content, err := getTemplateContent("")
		if assert.NoError(t, err) {
			assert.Equal(t, assets.OtelTmpl(), content)
		}
	})

	t.Run("from_file", func(t *testing.T) {
		t.Parallel()

		content, err := getTemplateContent("app_unit_test.go")
		if assert.NoError(t, err) && assert.NotEmpty(t, content) {
			assert.NotEqual(t, assets.OtelTmpl(), content)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		_, err := getTemplateContent("./not_found")
		assert.Error(t, err)
	})
}
