package info

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
)

const noBuffalo = "Warning: It seems like it is not a buffalo app. (.buffalo.dev.yml not found)\n\n"

// New returns a generator that performs buffalo
// related rx checks
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	devFile := filepath.Join(opts.App.Root, ".buffalo.dev.yml")
	if _, err := os.Stat(devFile); err != nil {
		g.RunFn(func(r *genny.Runner) error {
			opts.Out.Header("Buffalo: Application Details")
			return opts.Out.WriteString(noBuffalo)
		})
		return g, nil
	}

	g.RunFn(appDetails(opts))

	path := filepath.Join(opts.App.Root, "config")
	if _, err := os.Stat(path); err == nil {
		configFS := os.DirFS(path)
		g.RunFn(configs(opts, configFS))
	}

	aFS := os.DirFS(opts.App.Root)
	g.RunFn(pkgChecks(opts, aFS))

	return g, nil
}
