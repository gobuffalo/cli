package fix

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gobuffalo/cli/internal/runtime"
	"github.com/gobuffalo/genny/v2"
)

func FixDocker(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		if !opts.App.WithDocker {
			return nil
		}

		fmt.Println("~~~ Upgrading Dockerfile ~~~")
		dk, err := r.FindFile("Dockerfile")
		if err != nil {
			return err
		}

		ex := regexp.MustCompile(`(v[0-9.][\S]+)`)
		lines := strings.Split(dk.String(), "\n")
		for i, l := range lines {
			if strings.HasPrefix(strings.ToLower(l), "from gobuffalo/buffalo") {
				l = ex.ReplaceAllString(l, runtime.Version)
				lines[i] = l
			}
		}
		return r.File(genny.NewFileS(dk.Name(), strings.Join(lines, "\n")))
	}
}