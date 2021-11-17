package resource

import (
	"embed"
	"io/fs"
	"text/template"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

//go:embed templates/* templates/core/templates/folder-name/_form.plush.html.tmpl
var templates embed.FS

// New resource generator
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	if !opts.SkipTemplates {
		core, err := fs.Sub(templates, "templates/core")
		if err != nil {
			return g, err
		}

		if err := g.FS(core); err != nil {
			return g, err
		}
	}

	sub := "templates/use_model"
	if opts.SkipModel {
		sub = "templates/standard"
	}

	aFS, err := fs.Sub(templates, sub)
	if err != nil {
		return g, err
	}

	if err := g.FS(aFS); err != nil {
		return g, err
	}

	pres := presenter{
		App:   opts.App,
		Name:  name.New(opts.Name),
		Model: name.New(opts.Model),
		Attrs: opts.Attrs,
	}
	x := pres.Name.Resource().File().String()
	folder := pres.Name.Folder().Pluralize().String()
	g.Transformer(genny.Replace("resource-name", x))
	g.Transformer(genny.Replace("resource-use_model", x))
	g.Transformer(genny.Replace("folder-name", folder))

	data := map[string]interface{}{
		"opts":    pres,
		"actions": actions(opts),
		"folder":  folder,
	}
	helpers := template.FuncMap{
		"camelize": func(s string) string {
			return flect.Camelize(s)
		},
	}
	g.Transformer(gogen.TemplateTransformer(data, helpers))

	g.RunFn(installPop(opts))

	g.RunFn(addResource(pres))
	return g, nil
}
