package web

import (
	"embed"
	"html/template"
	"io/fs"

	"github.com/gobuffalo/cli/internal/genny/assets/standard"
	"github.com/gobuffalo/cli/internal/genny/assets/webpack"
	"github.com/gobuffalo/cli/internal/genny/newapp/core"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

//go:embed templates/* templates/templates/_flash.plush.html.tmpl
var templates embed.FS

// New generator for creating a Buffalo Web application
func New(opts *Options) (*genny.Group, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	gg, err := core.New(opts.Options)
	if err != nil {
		return gg, err
	}

	g := genny.New()
	g.Transformer(genny.Dot())
	data := map[string]interface{}{
		"opts": opts,
	}

	helpers := template.FuncMap{}

	t := gogen.TemplateTransformer(data, helpers)
	g.Transformer(t)
	sub, err := fs.Sub(templates, "templates")
	if err != nil {
		return gg, err
	}

	if err := g.FS(sub); err != nil {
		return gg, err
	}

	gg.Add(g)

	if opts.Webpack != nil {
		// add the webpack generator
		g, err = webpack.New(opts.Webpack)
		if err != nil {
			return gg, err
		}
		gg.Add(g)
	}

	if opts.Standard != nil {
		// add the standard generator
		g, err = standard.New(opts.Standard)
		if err != nil {
			return gg, err
		}
		gg.Add(g)
	}

	return gg, nil
}
