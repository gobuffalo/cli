package info

import (
	"io"
	"io/fs"

	"github.com/gobuffalo/genny/v2"
)

func pkgChecks(opts *Options, dir fs.FS) genny.RunFn {
	return func(r *genny.Runner) error {
		for _, x := range []string{"go.mod"} {
			f, err := dir.Open(x)
			if err != nil {
				return nil
			}
			s, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			opts.Out.Header("\nBuffalo: " + x)
			_, err = opts.Out.Write(s)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
