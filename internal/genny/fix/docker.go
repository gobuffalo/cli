package fix

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/gobuffalo/cli/internal/genny/docker"
	"github.com/gobuffalo/genny/v2"
)

func FixDocker(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		if !opts.App.WithDocker {
			return nil
		}

		fmt.Println("~~~ Checking Dockerfile ~~~")

		templates, err := docker.Templates()
		if err != nil {
			return err
		}

		tmpl, err := template.New("Dockerfile").ParseFS(templates, "Dockerfile.tmpl")
		if err != nil {
			return err
		}

		dockerOpts := &docker.Options{
			App: opts.App,
		}
		if err := dockerOpts.Validate(); err != nil {
			return err
		}

		bb := &bytes.Buffer{}
		if err := tmpl.ExecuteTemplate(bb, "Dockerfile.tmpl", map[string]interface{}{
			"opts": dockerOpts,
		}); err != nil {
			return err
		}

		f, err := r.FindFile("Dockerfile")
		if err != nil {
			return nil
		}

		if string(f.String()) == bb.String() {
			return nil
		}

		if !opts.YesToAll && !ask("Your Dockerfile is different from the latest Buffalo template.\nWould you like to replace yours with the latest template?") {
			fmt.Println("\tSkipping Dockerfile")
			return nil
		}

		_, err = f.Write(bb.Bytes())
		if err != nil {
			return err
		}
		return r.File(f)
	}
}
