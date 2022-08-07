package app_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hedhyw/otelinj/internal/app"
	"github.com/hedhyw/otelinj/internal/features"
	"github.com/hedhyw/otelinj/internal/pkg/config"
)

func TestAppContextDefined(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
package example

func Run(ctx context.Context) error {}
		`,
		Expected: `
package example
					
import (
	"github.com/hedhyw/otelinj/pkg/otelinji"
	"go.opentelemetry.io/otel"
)

func Run(ctx context.Context) error {
	ctx, span := otel.Tracer("example").Start(ctx, "Run")
	defer span.End()

	_ = ctx

}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppDifferentContextName(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
package example

func Run(context context.Context) error {}
		`,
		Expected: `
package example

import (
	"github.com/hedhyw/otelinj/pkg/otelinji"
	"go.opentelemetry.io/otel"
)

func Run(context context.Context) error {
	context, span := otel.Tracer("example").Start(context, "Run")
	defer span.End()

	_ = context

}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppSkippedName(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
package example

func Run(_ context.Context) error {}
		`,
		Expected: `
package example

import (
	"github.com/hedhyw/otelinj/pkg/otelinji"
	"go.opentelemetry.io/otel"
)

func Run(ctx context.Context) error {
	ctx, span := otel.Tracer("example").Start(ctx, "Run")
	defer span.End()

	_ = ctx

}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppReceiverNameValue(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
package example

type core struct {}

func (core) Run(ctx context.Context) error {}
		`,
		Expected: `
package example

import (
	"github.com/hedhyw/otelinj/pkg/otelinji"
	"go.opentelemetry.io/otel"
)

type core struct{}

func (core) Run(ctx context.Context) error {
	ctx, span := otel.Tracer("example").Start(ctx, "core.Run")
	defer span.End()

	_ = ctx

}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppReceiverNamePointer(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
package example

type core struct {}

func (*core) Run(ctx context.Context) error {}
		`,
		Expected: `
package example

import (
	"github.com/hedhyw/otelinj/pkg/otelinji"
	"go.opentelemetry.io/otel"
)

type core struct{}

func (*core) Run(ctx context.Context) error {
	ctx, span := otel.Tracer("example").Start(ctx, "core.Run")
	defer span.End()

	_ = ctx

}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppNoContextName(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
package example

func Run(context.Context) error {}
		`,
		Expected: `
package example

import (
	"github.com/hedhyw/otelinj/pkg/otelinji"
	"go.opentelemetry.io/otel"
)

func Run(ctx context.Context) error {
	ctx, span := otel.Tracer("example").Start(ctx, "Run")
	defer span.End()

	_ = ctx

}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppCtxUsed(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
package example

import "context"

func Run(ctx context.Context) error {
	ctx,cancel := context.WithCancel(ctx)
	cancel()
}
		`,
		Expected: `
package example

import (
	"context"

	"github.com/hedhyw/otelinj/pkg/otelinji"
	"go.opentelemetry.io/otel"
)

func Run(ctx context.Context) error {
	ctx, span := otel.Tracer("example").Start(ctx, "Run")
	defer span.End()

	ctx, cancel := context.WithCancel(ctx)
	cancel()
}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppErrParam(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
package example

func Run(ctx context.Context) (err error) {}
		`,
		Expected: `
package example

import (
	"github.com/hedhyw/otelinj/pkg/otelinji"
	"go.opentelemetry.io/otel"
)

func Run(ctx context.Context) (err error) {
	ctx, span := otel.Tracer("example").Start(ctx, "Run")
	defer func() { otelinji.EndSpanWithErr(span, err) }()

	_ = ctx

}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppErrParamDifferentName(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
package example

func Run(ctx context.Context) (err0 error) {}
		`,
		Expected: `
package example

import (
	"github.com/hedhyw/otelinj/pkg/otelinji"
	"go.opentelemetry.io/otel"
)

func Run(ctx context.Context) (err0 error) {
	ctx, span := otel.Tracer("example").Start(ctx, "Run")
	defer func() { otelinji.EndSpanWithErr(span, err0) }()

	_ = ctx

}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppComment(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
package example

// Example comment.
func Run(ctx context.Context) (error) {
	// Example comment.
}
		`,
		Expected: `
package example

import (
	"github.com/hedhyw/otelinj/pkg/otelinji"
	"go.opentelemetry.io/otel"
)

// Example comment.
func Run(ctx context.Context) error {
	// Example comment.
	ctx, span := otel.Tracer("example").Start(ctx, "Run")
	defer span.End()

	_ = ctx

}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppAlreadyImported(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
package example

import (
	"context"

	"github.com/hedhyw/otelinj/pkg/otelinji"

	"go.opentelemetry.io/otel"
)

func Run(ctx context.Context) (error) {}
		`,
		Expected: `
package example

import (
	"context"

	"github.com/hedhyw/otelinj/pkg/otelinji"

	"go.opentelemetry.io/otel"
)

func Run(ctx context.Context) error {
	ctx, span := otel.Tracer("example").Start(ctx, "Run")
	defer span.End()

	_ = ctx

}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppSkipGenerated(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
// Code generated by test v0.0.1. DO NOT EDIT.
package example

func Run(ctx context.Context) (error) {}
		`,
		Expected: `
// Code generated by test v0.0.1. DO NOT EDIT.
package example

func Run(ctx context.Context) error {}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppProcessGenerated(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
// Code generated by test v0.0.1. DO NOT EDIT.
package example

func Run(ctx context.Context) (error) {}
		`,
		Expected: `
// Code generated by test v0.0.1. DO NOT EDIT.
package example

import (
	"github.com/hedhyw/otelinj/pkg/otelinji"
	"go.opentelemetry.io/otel"
)

func Run(ctx context.Context) error {
	ctx, span := otel.Tracer("example").Start(ctx, "Run")
	defer span.End()

	_ = ctx

}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: false,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppCImport(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
package example

import "C"

func Run(ctx context.Context) (error) {}
		`,
		Expected: `
package example

import (
	"github.com/hedhyw/otelinj/pkg/otelinji"
	"go.opentelemetry.io/otel"
)

import "C"

func Run(ctx context.Context) error {
	ctx, span := otel.Tracer("example").Start(ctx, "Run")
	defer span.End()

	_ = ctx

}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppCustomTemplate(t *testing.T) {
	t.Parallel()

	const tmpl = `
package template

import "fmt"

func main() {
	fmt.Println("Hello, {{.FuncName}}")
}
`

	appTestCase{
		Input: `
package example

func World() () {}
		`,
		Expected: `
package example

import "fmt"

func World() {
	fmt.Println("Hello, World")
}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  features.RequireTempFile(t, []byte(tmpl)),
	})
}

func TestAppTwoFunctions(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
package example

func RunFirst(ctx context.Context) error {}

func RunSecond(ctx context.Context) error {}
		`,
		Expected: `
package example

import (
	"github.com/hedhyw/otelinj/pkg/otelinji"
	"go.opentelemetry.io/otel"
)

func RunFirst(ctx context.Context) error {
	ctx, span := otel.Tracer("example").Start(ctx, "RunFirst")
	defer span.End()

	_ = ctx

}

func RunSecond(ctx context.Context) error {
	ctx, span := otel.Tracer("example").Start(ctx, "RunSecond")
	defer span.End()

	_ = ctx

}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

func TestAppNothingChanged(t *testing.T) {
	t.Parallel()

	appTestCase{
		Input: `
package main

func main() {}
		`,
		Expected: `
package main

func main() {}
		`,
	}.Run(t, &config.Config{
		WriteIntoFile: false,
		SkipGenerated: true,
		FileName:      "",
		TemplateName:  "",
	})
}

type appTestCase struct {
	Input    string
	Expected string
}

// nolint: thelper // Not a helper.
func (tc appTestCase) Run(tb testing.TB, cfg *config.Config) {
	cfg.FileName = features.RequireTempFile(tb, []byte(tc.Input))

	var out bytes.Buffer

	err := app.New(cfg).Run(&out)
	require.NoError(tb, err)

	expected := features.RequireFormatGoSource(tb, tc.Expected)
	actual := features.RequireFormatGoSource(tb, out.String())

	assert.Equal(tb, expected, actual, actual)
}
