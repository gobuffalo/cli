package refresh

import (
	"embed"
	"strings"

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

	// TODO: workaround for 1.16, remove when we upgrade to 1.17 and rename "dot-*" files back to "-dot-*"
	g.Transformer(genny.NewTransformer("*", func(f genny.File) (genny.File, error) {
		name := f.Name()
		if strings.HasPrefix(name, "dot-") {
			name = strings.TrimPrefix(name, "dot-")
			name = "." + name
		}
		return genny.NewFile(name, f), nil
	}))
	g.Transformer(genny.Replace("/dot-", "/."))

	return g, nil
}
