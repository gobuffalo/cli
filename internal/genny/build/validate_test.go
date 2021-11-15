package build

import (
	"testing"

	"github.com/gobuffalo/cli/internal/genny/build/_fixtures/testtemplate"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/stretchr/testify/require"
)

func Test_TemplateValidator_Good(t *testing.T) {
	r := require.New(t)

	tvs := []TemplateValidator{PlushValidator}

	run := gentest.NewRunner()
	run.WithRun(ValidateTemplates(testtemplate.Good(), tvs))
	r.NoError(run.Run())
}

func Test_TemplateValidator_Bad(t *testing.T) {
	r := require.New(t)

	tvs := []TemplateValidator{PlushValidator}

	run := gentest.NewRunner()
	run.WithRun(ValidateTemplates(testtemplate.Bad(), tvs))

	err := run.Run()
	r.Error(err)
	r.Equal("template error in file a.html: line 1: no prefix parse function for > found\ntemplate error in file b.md: line 1: no prefix parse function for > found", err.Error())
}
