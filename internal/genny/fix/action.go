package fix

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

// ReplaceAppOnce fixes `actions/app.go` to fix the double execution issue.
// it covers https://github.com/gobuffalo/buffalo/issues/1653 and
// https://github.com/gobuffalo/cli/issues/228
func ReplaceAppOnce(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		fmt.Println("~~~ Apply AppOnce ~~~")

		file, err := r.FindFile("actions/app.go")
		if err != nil {
			return err
		}

		if !strings.Contains(file.String(), "appOnce.Do") {
			file, err = gogen.ReplaceBlock(file, "if app == nil {", "}", "appOnce.Do(func() {", "})")
			if err != nil {
				return err
			}

			file, err = gogen.AddGlobal(file, "appOnce sync.Once")
			if err != nil {
				return err
			}

			file, err = gogen.AddImport(file, "sync")
			if err != nil {
				return err
			}
		}

		return r.File(file)
	}
}
