package vcs

import (
	"strings"
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/stretchr/testify/require"
)

func Test_New_Bzr(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Provider: "bzr",
	})
	r.NoError(err)

	run := gentest.NewRunner()
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Equal(".bzrignore", f.Name())

	cmds := []string{
		"bzr init",
		"bzr add . -q",
		"bzr commit -q -m Initial Commit",
	}

	r.Len(res.Commands, len(cmds))
	for i, c := range res.Commands {
		r.Equal(cmds[i], strings.Join(c.Args, " "))
	}
}

func Test_New_Git(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Provider: "git",
	})
	r.NoError(err)

	run := gentest.NewRunner()
	r.NoError(run.With(g))
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Equal(".gitignore", f.Name())

	cmds := []string{
		"git init",
		"git add .",
		"git commit -q -m Initial Commit",
	}

	r.Len(res.Commands, len(cmds))
	for i, c := range res.Commands {
		r.Equal(cmds[i], strings.Join(c.Args, " "))
	}
}
