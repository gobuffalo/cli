package api

import (
	"gobuffalo/cli/internal/genny/newapp/core"
)

// Options for API applications
type Options struct {
	*core.Options
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.Options == nil {
		opts.Options = &core.Options{}
	}
	return opts.Options.Validate()
}
