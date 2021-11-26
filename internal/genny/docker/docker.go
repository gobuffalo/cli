package docker

import (
	"embed"
	"io/fs"
	"strings"
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

	sub, err := fs.Sub(templates, "templates")
	if err != nil {
		return g, err
	}

	if err := g.FS(sub); err != nil {
		return g, err
	}

	helpers := template.FuncMap{}
	t := gogen.TemplateTransformer(data, helpers)

	g.Transformer(t)
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
