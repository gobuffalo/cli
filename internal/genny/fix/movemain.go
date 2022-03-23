package fix

import (
	"fmt"

	"github.com/gobuffalo/genny/v2"
)

// MoveMain will move the main.go from the root folder into
// cmd/app/main.go
func MoveMain(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		// If there is a main in the cmd/app folder, we don't need to do anything
		_, err := r.FindFile("cmd/app/main.go")
		if err == nil {
			return nil
		}

		fmt.Println("~~~ Moving main.go ~~~")
		// If there is a main in the root folder, we need to move it
		f, err := r.FindFile("main.go")
		if err != nil {
			// There is no file to move to move\
			r.Logger.Info("No main.go found")
			return nil
		}

		nf := genny.NewFileS("cmd/app/main.go", f.String())
		r.Disk.Remove("main.go")

		return r.File(nf)
	}
}
