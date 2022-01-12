package fix

import (
	"bytes"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/gobuffalo/genny/v2"
)

func EncodeAppToml(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		p := "config/buffalo-app.toml"
		if _, err := r.FindFile(p); err == nil {
			return nil
		}
		dir := genny.NewDir(filepath.Dir(p), 0o755)
		if err := r.File(dir); err != nil {
			return err
		}

		bb := &bytes.Buffer{}
		if err := toml.NewEncoder(bb).Encode(opts.App); err != nil {
			return err
		}

		f := genny.NewFile(p, bb)
		return r.File(f)
	}
}
