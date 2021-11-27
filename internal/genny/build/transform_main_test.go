package build

import (
	"testing"

	"github.com/gobuffalo/cli/internal/genny/testrunner"
	"github.com/stretchr/testify/require"
)

func Test_transformMain(t *testing.T) {
	r := require.New(t)

	run, err := testrunner.WebApp()
	r.NoError(err)

	ref, err := testrunner.WebApp()
	main, err := ref.Disk.Find("main.go")
	r.NoError(err)
	run.Disk.Add(main)

	opts := &Options{}
	run.WithRun(transformMain(opts))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Files, 1)
	f := res.Files[0]
	r.Contains(f.String(), "func originalMain()")
}
