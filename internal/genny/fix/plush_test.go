package fix

import (
	"testing"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_UpdatePlushTemplates_AddsExtensions(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		App: meta.App{
			Root: ".",
		},
	}

	run := gentest.NewRunner()
	r.NoError(run.File(genny.NewFileS("templates/test1.html", "")))
	r.NoError(run.File(genny.NewFileS("templates/test2.js", "")))
	r.NoError(run.File(genny.NewFileS("templates/test3.md", "")))

	run.WithRun(UpdatePlushTemplates(opts))
	r.NoError(run.Run())

	results := run.Results()

	r.Equal("templates/test1.plush.html", results.Files[0].Name())
	r.Equal("templates/test2.plush.js", results.Files[1].Name())
	r.Equal("templates/test3.plush.md", results.Files[2].Name())
}

func Test_UpdatePlushTemplates_IgnoresNonWebFiles(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		App: meta.App{
			Root: ".",
		},
	}

	run := gentest.NewRunner()
	r.NoError(run.File(genny.NewFileS("templates/test1.go", "")))
	r.NoError(run.File(genny.NewFileS("templates/test2.txt", "")))

	run.WithRun(UpdatePlushTemplates(opts))
	r.NoError(run.Run())

	results := run.Results()

	r.Equal("templates/test1.go", results.Files[0].Name())
	r.Equal("templates/test2.txt", results.Files[1].Name())
}

func Test_UpdatePlushTemplates_IgnoresFilesWithFizzOrPlushExtension(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		App: meta.App{
			Root: ".",
		},
	}

	run := gentest.NewRunner()
	r.NoError(run.File(genny.NewFileS("templates/test1.fizz", "")))
	r.NoError(run.File(genny.NewFileS("templates/test2.plush.html", "")))

	run.WithRun(UpdatePlushTemplates(opts))
	r.NoError(run.Run())

	results := run.Results()

	r.Equal("templates/test1.fizz", results.Files[0].Name())
	r.Equal("templates/test2.plush.html", results.Files[1].Name())
}