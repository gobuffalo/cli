package fix

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/gobuffalo/cli/internal/genny/newapp/api"
	"github.com/gobuffalo/cli/internal/genny/newapp/web"
	"github.com/gobuffalo/genny/v2"
)

func FixEmbed(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		fmt.Println("~~~ Checking embed.go files ~~~")

		files := []string{"locales/embed.go", "public/embed.go", "templates/embed.go"}
		templates, err := web.Templates()
		if err != nil {
			return err
		}

		if opts.App.AsAPI {
			files = []string{"locales/embed.go"}

			templates, err = api.Templates()
			if err != nil {
				return err
			}
		} else {
			f := genny.NewFileS("templates/home/delete_me.txt", "you can delete this file")
			if err := r.File(f); err != nil {
				return err
			}
		}

		for _, name := range files {
			tmpl, err := template.New("embed.go").ParseFS(templates, name+".tmpl")
			if err != nil {
				return err
			}

			bb := &bytes.Buffer{}
			if err := tmpl.ExecuteTemplate(bb, "embed.go.tmpl", nil); err != nil {
				return err
			}

			f := genny.NewFile(name, bb)
			if err := r.File(f); err != nil {
				return err
			}
		}
		return nil
	}
}
