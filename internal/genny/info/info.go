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

	configFS := os.DirFS(filepath.Join(opts.App.Root, "config"))
	g.RunFn(configs(opts, configFS))

	aFS := os.DirFS(opts.App.Root)
	g.RunFn(pkgChecks(opts, aFS))

	return g, nil
}
