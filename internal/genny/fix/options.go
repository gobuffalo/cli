package fix

import (
	"github.com/gobuffalo/meta"
)

// Options for building a Buffalo application
type Options struct {
	App meta.App `json:"app"`
	// YesToAll will be used by the command to skip the confirmation
	// and perform all implied destroy operations
	YesToAll bool `json:"yes_to_all"`
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}
	return nil
}
