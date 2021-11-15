package info

import (
	"io/fs"
	"path"

	"github.com/gobuffalo/genny/v2"
)

func configs(opts *Options, dir fs.FS) genny.RunFn {
	return func(r *genny.Runner) error {
		return fs.WalkDir(dir, ".", func(p string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if entry.IsDir() {
				return nil
			}

			b, err := fs.ReadFile(dir, p)
			if err != nil {
				return err
			}
			opts.Out.Header("Buffalo: " + path.Join("config", p))
			_, err = opts.Out.Write(append(b, '\n'))
			return err
		})
	}
}
