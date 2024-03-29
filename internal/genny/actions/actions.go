package actions

import (
	"embed"
	"fmt"
	"strings"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

//go:embed templates/*
var templates embed.FS

// New returns a new generator for build actions on a Buffalo app
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	g.RunFn(construct(opts))
	return g, nil
}

func construct(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		pres := &presenter{
			Name:    name.New(opts.Name),
			Data:    data{},
			Helpers: data{},
			Options: opts,
		}

		if err := buildActions(pres)(r); err != nil {
			return err
		}

		if err := buildTests(pres)(r); err != nil {
			return err
		}

		if err := updateApp(pres)(r); err != nil {
			return err
		}

		if !opts.SkipTemplates {
			if err := buildTemplates(pres)(r); err != nil {
				return err
			}
		}
		return nil
	}
}

func transform(pres *presenter, f genny.File) (genny.File, error) {
	pres.Data["actions"] = pres.Actions
	pres.Data["name"] = pres.Name
	t := gogen.TemplateTransformer(pres.Data, pres.Helpers)
	return t.Transform(f)
}

func updateApp(pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		f, err := r.FindFile("actions/app.go")
		if err != nil {
			return err
		}

		var lines []string
		body := f.String()
		for _, a := range pres.Actions {
			e := fmt.Sprintf("app.%s(\"/%s/%s\", %s%s)", strings.ToUpper(pres.Options.Method), pres.Name.Underscore(), a.Underscore(), pres.Name.Pascalize(), a.Pascalize())
			if !strings.Contains(body, e) {
				lines = append(lines, e)
			}
		}

		f, err = gogen.AddInsideBlock(f, "appOnce.Do(func() {", strings.Join(lines, "\n\t\t"))
		if err != nil {
			if strings.Contains(err.Error(), "could not find desired block") {
				f, err = gogen.AddInsideBlock(f, "if app == nil {", strings.Join(lines, "\n\t\t"))
				if err != nil {
					return err
				} else {
					r.Logger.Warnf("This app was built with CLI v0.18.8 or older. See https://gobuffalo.io/documentation/known-issues/#cli-v0.18.8")
				}
			} else {
				return err
			}
		}
		return r.File(f)
	}
}
