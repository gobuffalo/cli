package docker

import (
	"embed"
	"io/fs"
	"text/template"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

//go:embed templates/*
var templates embed.FS

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	data := map[string]interface{}{
		"opts": opts,
	}

	common, err := fs.Sub(templates, "templates/common")
	if err != nil {
		return g, err
	}

	if err := g.FS(common); err != nil {
		return g, err
	}

	switch opts.Style {
	case "multi":
		multi, err := fs.Sub(templates, "templates/multi")
		if err != nil {
			return g, err
		}

		if err := g.FS(multi); err != nil {
			return g, err
		}
	case "standard":
		standard, err := fs.Sub(templates, "templates/standard")
		if err != nil {
			return g, err
		}

		if err := g.FS(standard); err != nil {
			return g, err
		}
	}

	helpers := template.FuncMap{}
	t := gogen.TemplateTransformer(data, helpers)
	g.Transformer(t)
	g.Transformer(genny.Dot())
	return g, nil
}
