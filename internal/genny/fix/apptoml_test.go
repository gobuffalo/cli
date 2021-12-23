package fix

import (
	"testing"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_EncodeAppToml(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()

	g := EncodeAppToml(&Options{
		App: meta.App{
			Name: name.New("coke"),
		},
	})
	run.WithRun(g)

	r.NoError(run.Run())

	results := run.Results()
	f, err := results.Find("config/buffalo-app.toml")
	r.NoError(err)

	r.Contains(f.String(), `name = "coke"`)
}

func Test_EncodeAppToml_NotReplaceExisting(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()
	run.Disk.Add(genny.NewFileS("config/buffalo-app.toml", `name = "pepsi"`))

	g := EncodeAppToml(&Options{
		App: meta.App{
			Name: name.New("coke"),
		},
	})
	run.WithRun(g)

	r.NoError(run.Run())

	results := run.Results()
	f, err := results.Find("config/buffalo-app.toml")
	r.NoError(err)

	r.Contains(f.String(), `name = "pepsi"`)
}
