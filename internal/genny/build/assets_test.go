package build

import (
	"archive/zip"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/gobuffalo/cli/internal/genny/assets/webpack"
	"github.com/gobuffalo/cli/internal/genny/newapp/core"
	"github.com/gobuffalo/cli/internal/genny/newapp/web"
	"github.com/gobuffalo/cli/internal/genny/testrunner"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_assets(t *testing.T) {
	r := require.New(t)

	run, err := testrunner.WebApp(&web.Options{})
	r.NoError(err)

	opts := &Options{
		WithAssets: true,
	}
	r.NoError(opts.Validate())
	opts.App.WithNodeJs = true
	opts.App.PackageJSON.Scripts = map[string]string{
		"build": "webpack -p --progress",
	}
	r.NoError(run.WithNew(assets(opts)))

	r.NoError(envy.MustSet("NODE_ENV", ""))
	ne := envy.Get("NODE_ENV", "")
	r.Empty(ne)
	r.NoError(run.Run())

	ne = envy.Get("NODE_ENV", "")
	r.NotEmpty(ne)
	r.Equal(opts.Environment, ne)

	res := run.Results()

	cmds := []string{"npm run build"}
	r.Len(res.Commands, len(cmds))
	for i, c := range res.Commands {
		r.Equal(cmds[i], strings.Join(c.Args, " "))
	}
}

func Test_assets_Archived(t *testing.T) {
	r := require.New(t)

	app := meta.New(".")
	app.Root = "."
	webOpts := &web.Options{
		Options: &core.Options{},
		Webpack: &webpack.Options{},
	}
	webOpts.App = app
	run, err := testrunner.WebApp(webOpts)
	r.NoError(err)

	opts := &Options{
		WithAssets:    true,
		ExtractAssets: true,
	}
	opts.App = app
	r.NoError(opts.Validate())
	r.NoError(run.WithNew(assets(opts)))
	r.NoError(run.Run())

	res := run.Results()

	cmds := []string{}
	r.Len(res.Commands, len(cmds))
	for i, c := range res.Commands {
		r.Equal(cmds[i], strings.Join(c.Args, " "))
	}

	f, err := res.Find("actions/app.go")
	r.NoError(err)
	r.Contains(f.String(), `// app.ServeFiles("/"`)

	f, err = res.Find("bin/assets.zip")
	r.NoError(err)

	tmp, err := os.CreateTemp("", "assets-*.zip")
	r.NoError(err)
	t.Cleanup(func() {
		os.Remove(tmp.Name())
	})

	_, err = io.Copy(tmp, f)
	r.NoError(err)
	r.NoError(tmp.Close())

	archive, err := zip.OpenReader(tmp.Name())
	r.NoError(err)

	r.Equal(1, len(archive.File))
	for _, e := range []string{"keep"} {
		_, err = archive.Open(e)
		r.NoError(err)
	}
}
