package add

import (
	"bytes"

	"github.com/gobuffalo/cli/internal/plugins/plugdeps"
	"github.com/gobuffalo/genny/v2"
)

// New add plugin to the config file
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	bb := &bytes.Buffer{}
	plugs := plugdeps.New()
	plugs.Add(opts.Plugins...)
	if err := plugs.Encode(bb); err != nil {
		return g, err
	}

	g.File(genny.NewFile(plugdeps.ConfigPath(opts.App), bb))
	return g, nil
}
