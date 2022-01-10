package fix

import (
	"testing"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_WebpackCheck_NoOverwriteExisting(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		App: meta.App{
			Root:        ".",
			WithWebpack: true,
		},
	}

	run := gentest.NewRunner()
	bb, err := defaultWebpack(opts.App)
	r.NoError(err)
	fileContents := bb.String()
	r.NoError(run.File(genny.NewFileS("webpack.config.js", fileContents)))

	run.WithRun(WebpackCheck(opts))
	r.NoError(run.Run())

	results := run.Results()
	f := results.Files[0]
	r.Equal("webpack.config.js", f.Name())
	r.Equal(fileContents, f.String())
}

func Test_WebpackCheck_UpdatingFile(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		App: meta.App{
			Root:        ".",
			WithWebpack: true,
		},
		YesToAll: true,
	}

	bb, err := defaultWebpack(opts.App)
	r.NoError(err)
	fileContents := bb.String()

	run := gentest.NewRunner()
	r.NoError(run.File(genny.NewFileS("webpack.config.js", "console.log('hello world')")))

	run.WithRun(WebpackCheck(opts))
	r.NoError(run.Run())

	results := run.Results()

	f := results.Files[0]
	r.Equal("webpack.config.js", f.Name())
	r.Equal(fileContents, f.String())
}
