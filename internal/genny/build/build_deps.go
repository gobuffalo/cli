package build

import (
	"github.com/gobuffalo/genny/v2"
)

func buildDeps(opts *Options) (*genny.Generator, error) {
	g := genny.New()
	err := opts.Validate()

	return g, err
}
