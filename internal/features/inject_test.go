package features_test

import (
	_ "embed"
	"testing"

	"github.com/hedhyw/gherkingen/v2/pkg/bdd"
	"github.com/stretchr/testify/assert"

	"github.com/hedhyw/otelinji/internal/features"
)

//go:embed inject.exp.test
var expectedInjectSource string

func TestAddingOpentracingDefinition(t *testing.T) {
	t.Parallel()

	f := bdd.NewFeature(t, "Adding opentracing definition")

	f.Scenario("User wants to add opentelemetry layers", func(t *testing.T, f *bdd.Feature) {
		t.Parallel()

		var args []string
		var output string

		f.When("the user provides an input file `inject.in.test`", func() {
			args = append(args, "-filename", "./inject.in.test")
		})
		f.And("runs the application", func() {
			output = features.RequireRunApp(t, args...)
		})
		f.Then("the output equals to the content of the file `inject.exp.test`", func() {
			expected := features.RequireFormatGoSource(t, expectedInjectSource)
			actual := features.RequireFormatGoSource(t, output)

			assert.Equal(t, expected, actual, actual)
		})
	})
}
