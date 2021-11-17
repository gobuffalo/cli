package refresh

import (
	"embed"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/plushgen"
	"github.com/gobuffalo/plush/v4"
)

//go:embed templates/*
var templates embed.FS

// New generator to generate refresh templates
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()
	if err := opts.Validate(); err != nil {
		return g, err
	}

	if err := g.FS(templates); err != nil {
		return g, err
	}

	ctx := plush.NewContext()
	ctx.Set("app", opts.App)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Dot())
	return g, nil
}
