package features_test

import (
	_ "embed"
	"testing"

	"github.com/hedhyw/gherkingen/v2/pkg/bdd"
	"github.com/stretchr/testify/assert"

	"github.com/hedhyw/otelinji/internal/features"
)

//go:embed inject.exp.go.txt
var expectedInjectSource string

func TestAddingOpentracingDefinitionToTheFile(t *testing.T) {
	t.Parallel()

	f := bdd.NewFeature(t, "Adding opentracing definition to the file")

	f.Scenario("User wants to add opentelemetry layers", func(t *testing.T, f *bdd.Feature) {
		t.Parallel()

		var (
			args   []string
			output string
		)

		f.Given("the user provides an input file `inject.exp.go.txt`", func() {
			args = append(args, "-filename", "./inject.in.go.txt")
		})
		f.When("he runs the application", func() {
			output = features.RequireRunApp(t, args...)
		})
		f.Then("the output equals to the content of the file `inject.in.go.txt`", func() {
			expected := features.RequireFormatGoSource(t, expectedInjectSource)
			actual := features.RequireFormatGoSource(t, output)

			assert.Equal(t, expected, actual, actual)
		})
	})
}
