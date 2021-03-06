package standard

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

// New generator for creating basic asset files
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	sub, err := fs.Sub(templates, "templates")
	if err != nil {
		return g, err
	}

	if err := g.FS(sub); err != nil {
		return g, err
	}

	data := map[string]interface{}{}
	h := template.FuncMap{}
	t := gogen.TemplateTransformer(data, h)
	g.Transformer(t)

	g.RunFn(func(r *genny.Runner) error {
		f, err := r.FindFile("templates/application.plush.html")
		if err != nil {
			return err
		}

		s := strings.Replace(f.String(), "</title>", "</title>\n"+bs4, 1)
		return r.File(genny.NewFileS(f.Name(), s))
	})

	return g, nil
}

const bs4 = `<link href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">`
