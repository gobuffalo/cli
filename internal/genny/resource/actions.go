package resource

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

func addResource(pres presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		f, err := r.FindFile("actions/app.go")
		if err != nil {
			return err
		}

		stmt := fmt.Sprintf("app.Resource(\"/%s\", %sResource{})", pres.Name.URL(), pres.Name.Resource())
		f, err = gogen.AddInsideBlock(f, "appOnce.Do(func() {", stmt)
		if err != nil {
			if strings.Contains(err.Error(), "could not find desired block") {
				f, err = gogen.AddInsideBlock(f, "if app == nil {", stmt)
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

func actions(opts *Options) []name.Ident {
	actions := []name.Ident{
		name.New("list"),
		name.New("show"),
		name.New("create"),
		name.New("update"),
		name.New("destroy"),
	}

	if opts.App.AsWeb {
		actions = append(actions, name.New("new"), name.New("edit"))
	}

	return actions
}
