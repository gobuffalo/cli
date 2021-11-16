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

func Test_pkgChecks(t *testing.T) {
	r := require.New(t)

	bb := &bytes.Buffer{}

	run := gentest.NewRunner()

	opts := &Options{
		App: meta.New("."),
		Out: rx.NewWriter(bb),
	}

	run.WithRun(pkgChecks(opts, os.DirFS("../info/testtemplate/module")))
	r.NoError(run.Run())

	res := bb.String()
	r.Contains(res, "Buffalo: go.mod")
}
