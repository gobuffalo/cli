package build

import (
	"testing"

	"github.com/gobuffalo/cli/internal/genny/newapp/web"
	"github.com/gobuffalo/cli/internal/genny/testrunner"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/stretchr/testify/require"
)

func Test_transformMain(t *testing.T) {
	r := require.New(t)

	ref, err := testrunner.WebApp(&web.Options{})
	r.NoError(err)
	main, err := ref.Disk.Find("cmd/app/main.go")
	r.NoError(err)

	run := gentest.NewRunner()
	run.Disk.Add(main)

	opts := &Options{}
	run.WithRun(transformMain(opts))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Files, 1)
	f := res.Files[0]
	r.Contains(f.String(), "func originalMain()")
}
