package info

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
)

// New returns a generator that performs buffalo
// related rx checks
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	g.RunFn(appDetails(opts))

	path := filepath.Join(opts.App.Root, "config")
	if err := os.MkdirAll(path, 0755); err != nil {
		return g, err
	}
	configFS := os.DirFS(path)
	g.RunFn(configs(opts, configFS))

	aFS := os.DirFS(opts.App.Root)
	g.RunFn(pkgChecks(opts, aFS))

	return g, nil
}
