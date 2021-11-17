package build

import (
	"io/fs"
	"testing"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/stretchr/testify/require"
)

func Test_transformMain(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()
	b, err := fs.ReadFile(coke, "main.go")
	r.NoError(err)

	run.Disk.Add(genny.NewFileB("main.go", b))

	opts := &Options{}
	run.WithRun(transformMain(opts))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Files, 1)
	f := res.Files[0]
	r.Contains(f.String(), "func originalMain()")
}
