package info

import (
	"bytes"
	"os"
	"testing"

	"github.com/gobuffalo/clara/v2/genny/rx"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_configs(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()

	bb := &bytes.Buffer{}

	app := meta.New(".")
	opts := &Options{
		App: app,
		Out: rx.NewWriter(bb),
	}

	run.WithRun(configs(opts, os.DirFS("../info/templates/config")))
	r.NoError(run.Run())

	x := bb.String()
	r.Contains(x, "Buffalo: config/buffalo-app.toml\napp")
	r.Contains(x, "Buffalo: config/buffalo-plugins.toml\nplugins")
}
