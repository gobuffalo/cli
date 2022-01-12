package build

import (
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/psanford/memfs"
	"github.com/stretchr/testify/require"
)

func Test_TemplateValidator_Good(t *testing.T) {
	r := require.New(t)

	tvs := []TemplateValidator{PlushValidator}

	goodFS := memfs.New()
	r.NoError(goodFS.MkdirAll("_ignored", 0o755))
	r.NoError(goodFS.WriteFile("_ignored/c.html", []byte("c"), 0o644))
	r.NoError(goodFS.WriteFile("a.html", []byte("a"), 0o644))
	r.NoError(goodFS.WriteFile("b.html", []byte("b"), 0o644))

	run := gentest.NewRunner()
	run.WithRun(ValidateTemplates(goodFS, tvs))
	r.NoError(run.Run())
}

func Test_TemplateValidator_Bad(t *testing.T) {
	r := require.New(t)

	tvs := []TemplateValidator{PlushValidator}

	badFS := memfs.New()
	r.NoError(badFS.WriteFile("a.html", []byte("A Hello <%= broken!>%>>"), 0o644))
	r.NoError(badFS.WriteFile("b.md", []byte("B Hello <%= broken!>%>>"), 0o644))

	run := gentest.NewRunner()
	run.WithRun(ValidateTemplates(badFS, tvs))

	err := run.Run()
	r.Error(err)
	r.Equal("template error in file a.html: line 1: no prefix parse function for > found\ntemplate error in file b.md: line 1: no prefix parse function for > found", err.Error())
}
