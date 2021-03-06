package grift

import (
	"text/template"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

// New generator to create a grift task
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	data := map[string]interface{}{
		"opts": opts,
	}
	t := gogen.TemplateTransformer(data, template.FuncMap{})
	g.Transformer(t)

	g.RunFn(func(r *genny.Runner) error {
		return genFile(r, opts)
	})
	return g, nil
}

func genFile(r *genny.Runner, opts *Options) error {
	header := tmplHeader
	path := "grifts/" + opts.Name.File(".go.tmpl").String()
	if f, err := r.FindFile(path); err == nil {
		header = f.String()
	}
	f := genny.NewFileS(path, header+tmplBody)
	return r.File(f)
}
