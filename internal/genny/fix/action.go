package fix

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

// ReplaceOldImports walks all the .go files in an application
// It will then attempt to convert any old import paths to any new import paths
// used by this version Buffalo.
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
