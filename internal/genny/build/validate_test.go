package build

import (
	"os"
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/stretchr/testify/require"
)

func Test_TemplateValidator_Good(t *testing.T) {
	r := require.New(t)

	tvs := []TemplateValidator{PlushValidator}

	run := gentest.NewRunner()
	run.WithRun(ValidateTemplates(os.DirFS("../build/_fixtures/template_validator/good"), tvs))
	r.NoError(run.Run())
}

func Test_TemplateValidator_Bad(t *testing.T) {
	r := require.New(t)

	tvs := []TemplateValidator{PlushValidator}

	run := gentest.NewRunner()
	run.WithRun(ValidateTemplates(os.DirFS("../build/_fixtures/template_validator/bad"), tvs))

	err := run.Run()
	r.Error(err)
	r.Equal("template error in file a.html: line 1: no prefix parse function for > found\ntemplate error in file b.md: line 1: no prefix parse function for > found", err.Error())
}
