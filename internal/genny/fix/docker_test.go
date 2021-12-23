package fix

import (
	"testing"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_FixDocker_SkipsExisting(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()
	run.Disk.Add(genny.NewFileS("Dockerfile", "my custom Dockerfile"))

	opts := &Options{
		App: meta.App{
			Root:       ".",
			WithDocker: true,
		},
		YesToAll: true,
	}
	g := FixDocker(opts)
	run.WithRun(g)

	r.NoError(run.Run())
	results := run.Results()

	f, err := results.Find("Dockerfile")
	r.NoError(err)

	r.Contains(f.String(), "multi-stage Dockerfile")
}
