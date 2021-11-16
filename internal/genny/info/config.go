package info

import (
	"io/fs"
	"path"

	"github.com/gobuffalo/genny/v2"
)

func configs(opts *Options, fsys fs.FS) genny.RunFn {
	return func(r *genny.Runner) error {
		return fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			b, err := fs.ReadFile(fsys, p)
			if err != nil {
				return err
			}
			opts.Out.Header("Buffalo: " + path.Join("config", p))
			_, err = opts.Out.Write(append(b, '\n'))
			return err
		})
	}
}
