package features_test

import (
	"testing"

	"github.com/hedhyw/otelinji/internal/features"

	"github.com/hedhyw/gherkingen/v2/pkg/bdd"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Parallel()

	f := bdd.NewFeature(t, "Version")

	/* As a user I want to know the application version. */

	f.Scenario("User wants to know the version", func(t *testing.T, f *bdd.Feature) {
		var (
			args   []string
			output string
		)

		f.Given("the user provides the flag \"-version\"", func() {
			args = append(args, "-version")
		})
		f.When("they run the application", func() {
			output = features.RequireRunApp(t, args...)
		})
		f.Then("the output contains the version of the server", func() {
			assert.Contains(t, output, "github.com/hedhyw/otelinji@v0.0.1")
		})
	})
}
