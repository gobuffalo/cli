package actions

import (
	"fmt"
	"io"

	"github.com/gobuffalo/genny/v2"
)

// buildTests is the top level action/test builder
// it determines whether to build a new actions/foo_test.go file
// or append to an existing one
func buildTests(pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		fn := fmt.Sprintf("actions/%s_test.go", pres.Name.File())
		xf, err := r.FindFile(fn)
		if err != nil {
			return buildNewTests(fn, pres)(r)
		}
		return appendTests(xf, pres)(r)
	}
}

// buildNewTests builds a brand new actions/foo_test.go file
// and files it with tests
func buildNewTests(fn string, pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		h, err := templates.Open("templates/tests_header.go.tmpl")
		if err != nil {
			return err
		}
		a, err := templates.Open("templates/test.go.tmpl")
		if err != nil {
			return err
		}

		f := genny.NewFile(fn+".tmpl", io.MultiReader(h, a))
		f, err = transform(pres, f)
		if err != nil {
			return err
		}
		return r.File(f)
	}
}

// appendTests appends *only* tests that don't exist in
// actions/foo_test.go. if the test already exists it is not touched.
func appendTests(f genny.File, pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		a, err := templates.Open("templates/test.go.tmpl")
		if err != nil {
			return err
		}

		f := genny.NewFile(f.Name()+".tmpl", io.MultiReader(f, a))
		f, err = transform(pres, f)
		if err != nil {
			return err
		}
		return r.File(f)
	}
}
