package fix

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/gobuffalo/genny/v2"
)

func EncodeAppToml(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		p := filepath.Join("config", "buffalo-app.toml")
		if _, err := os.Stat(p); err == nil {
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		f, err := os.Create(p)
		if err != nil {
			return err
		}
		if err := toml.NewEncoder(f).Encode(opts.App); err != nil {
			return err
		}
		return nil
	}
}
