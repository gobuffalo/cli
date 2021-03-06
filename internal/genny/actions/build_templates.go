package actions

import (
	"fmt"
	"io/fs"

	"github.com/gobuffalo/genny/v2"
)

func buildTemplates(pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		sub, err := fs.Sub(templates, "templates")
		if err != nil {
			return err
		}

		f, err := fs.ReadFile(sub, "view.plush.html.tmpl")
		if err != nil {
			return err
		}
		for _, a := range pres.Actions {
			pres.Data["action"] = a
			fn := fmt.Sprintf("templates/%s/%s.plush.html.tmpl", pres.Name.Folder(), a.File())
			xf := genny.NewFileB(fn, f)
			xf, err = transform(pres, xf)
			if err != nil {
				return err
			}
			if err := r.File(xf); err != nil {
				return err
			}
		}
		return nil
	}
}
