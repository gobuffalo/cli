package install

import (
	"bytes"
	"go/build"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/cli/internal/plugins/plugdeps"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		App: meta.New("."),
		Plugins: []plugdeps.Plugin{
			{Binary: "buffalo-pop", GoGet: "github.com/gobuffalo/buffalo-pop/v3@latest", Tags: meta.BuildTags{"sqlite"}},
			{Binary: "buffalo-hello.rb", Local: "./plugins/buffalo-hello.rb"},
		},
	})
	r.NoError(err)

	run := gentest.NewRunner()
	c := build.Default
	run.Disk.Add(genny.NewFile(filepath.Join(c.GOPATH, "bin", "buffalo-pop"), &bytes.Buffer{}))
	run.FileFn = func(f genny.File) (genny.File, error) {
		bb := &bytes.Buffer{}
		if _, err := io.Copy(bb, f); err != nil {
			return f, err
		}
		return genny.NewFile(f.Name(), bb), nil
	}

	run.WithGroup(g)

	r.NoError(run.Run())

	res := run.Results()

	ecmds := []string{"go install -tags sqlite github.com/gobuffalo/buffalo-pop/v3@latest"}
	r.Len(res.Commands, len(ecmds))
	for i, c := range res.Commands {
		r.Equal(ecmds[i], strings.Join(c.Args, " "))
	}

	efiles := []string{"config/buffalo-plugins.toml", "bin/buffalo-pop"}
	r.Len(res.Files, len(efiles))
	for i, f := range res.Files {
		r.True(strings.HasSuffix(f.Name(), efiles[i]))
	}
}
