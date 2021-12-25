package fix

import (
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_FixEmbed_ApiProject(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()

	opts := &Options{
		App: meta.App{
			Root:  ".",
			AsAPI: true,
		},
	}
	g := FixEmbed(opts)
	run.WithRun(g)

	r.NoError(run.Run())
	results := run.Results()

	f, err := results.Find("locales/embed.go")
	r.NoError(err)
	r.Contains(f.String(), "buffalo.NewFS")
}

func Test_FixEmbed_WebProject(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()

	opts := &Options{
		App: meta.App{
			Root: ".",
		},
	}
	g := FixEmbed(opts)
	run.WithRun(g)

	r.NoError(run.Run())
	results := run.Results()

	f, err := results.Find("locales/embed.go")
	r.NoError(err)
	r.Contains(f.String(), "buffalo.NewFS")

	f, err = results.Find("public/embed.go")
	r.NoError(err)
	r.Contains(f.String(), "buffalo.NewFS")

	f, err = results.Find("templates/embed.go")
	r.NoError(err)
	r.Contains(f.String(), "buffalo.NewFS")
}
