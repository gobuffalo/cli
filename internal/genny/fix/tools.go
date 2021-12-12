package fix

import (
	"fmt"
	"os/exec"

	"github.com/gobuffalo/genny/v2"
)

// InstallTools installs required tools like the pop plugin
func InstallTools(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		fmt.Println("~~~ Installing required tools ~~~")

		if opts.App.WithPop {
			if err := r.Exec(exec.Command("go", "install", "github.com/gobuffalo/buffalo-pop/v3@latest")); err != nil {
				return err
			}
		}

		return nil
	}
}
