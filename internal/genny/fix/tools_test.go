package fix

import (
	"strings"
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_InstallTools_WithOutPop(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()

	run.WithRun(InstallTools(&Options{}))

	r.NoError(run.Run())

	results := run.Results()
	r.Len(results.Commands, 0)
}

func Test_InstallTools_WithPop(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()

	run.WithRun(InstallTools(&Options{
		App: meta.App{
			WithPop: true,
		},
	}))

	r.NoError(run.Run())

	results := run.Results()
	r.Len(results.Commands, 1)
	c := results.Commands[0]
	r.Equal("go install github.com/gobuffalo/buffalo-pop/v3@latest", strings.Join(c.Args, " "))
}
