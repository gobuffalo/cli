package build

import (
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/gobuffalo/cli/internal/genny/testrunner"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

var eq = func(r *require.Assertions, s string, c *exec.Cmd) {
	if runtime.GOOS == "windows" {
		s = strings.Replace(s, "bin/build", `bin\build.exe`, 1)
		s = strings.Replace(s, "bin/foo", `bin\foo.exe`, 1)
	}
	r.Equal(s, strings.Join(c.Args, " "))
}

func Test_New(t *testing.T) {
	r := require.New(t)

	run, err := testrunner.WebApp()
	r.NoError(err)

	opts := &Options{
		WithAssets:    true,
		WithBuildDeps: true,
		Environment:   "bar",
		App:           meta.New("."),
	}
	opts.App.Bin = "bin/foo"
	r.NoError(run.WithNew(New(opts)))
	run.Root = opts.App.Root
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Files, 0)

	cmds := []string{
		"go mod tidy",
		"go build -tags bar -o bin/foo",
		"go mod tidy",
	}
	r.Len(res.Commands, len(cmds))
	for i, c := range res.Commands {
		eq(r, cmds[i], c)
	}
}

func Test_NewWithoutBuildDeps(t *testing.T) {
	envy.Temp(func() {
		r := require.New(t)

		run, err := testrunner.WebApp()
		r.NoError(err)

		opts := &Options{
			WithAssets:    false,
			WithBuildDeps: false,
			Environment:   "bar",
			App:           meta.New("."),
		}
		opts.App.Bin = "bin/foo"
		r.NoError(run.WithNew(New(opts)))
		run.Root = opts.App.Root

		r.NoError(run.Run())

		res := run.Results()

		cmds := []string{"go mod tidy", "go build -tags bar -o bin/foo"}
		r.Len(res.Commands, len(cmds))
		for i, c := range res.Commands {
			eq(r, cmds[i], c)
		}
	})
}
