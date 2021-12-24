package fix

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/gobuffalo/cli/internal/genny/assets/webpack"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/meta"
)

// WebpackCheck will compare the current default Buffalo
// webpack.config.js against the applications webpack.config.js. If they are
// different you have the option to overwrite the existing webpack.config.js
// file with the new one.
func WebpackCheck(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		if !opts.App.WithWebpack {
			return nil
		}

		fmt.Println("~~~ Checking webpack.config.js ~~~")
		bb, err := defaultWebpack(opts.App)
		if err != nil {
			return err
		}

		f, err := r.FindFile("webpack.config.js")
		if err != nil {
			return err
		}

		if f.String() == bb.String() {
			return nil
		}

		if !opts.YesToAll && !ask("Your webpack.config.js file is different from the latest Buffalo template.\nWould you like to replace yours with the latest template?") {
			fmt.Println("\tSkipping webpack.config.js")
			return nil
		}

		_, err = f.Write(bb.Bytes())
		return err
	}
}

func defaultWebpack(app meta.App) (*bytes.Buffer, error) {
	templates, err := webpack.Templates()
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("webpack").ParseFS(templates, "webpack.config.js.tmpl")
	if err != nil {
		return nil, err
	}

	bb := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(bb, "webpack.config.js.tmpl", map[string]interface{}{
		"opts": &webpack.Options{
			App: app,
		},
	})

	return bb, err
}
