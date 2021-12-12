package fix

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/gobuffalo/cli/internal/genny/assets/webpack"
	"github.com/gobuffalo/genny/v2"
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

		templates, err := webpack.Templates()
		if err != nil {
			return err
		}

		tmpl, err := template.New("webpack").ParseFS(templates, "webpack.config.js.tmpl")
		if err != nil {
			return err
		}

		bb := &bytes.Buffer{}
		err = tmpl.ExecuteTemplate(bb, "webpack.config.js.tmpl", map[string]interface{}{
			"opts": &webpack.Options{
				App: opts.App,
			},
		})
		if err != nil {
			return err
		}

		b, err := os.ReadFile("webpack.config.js")
		if err != nil {
			return err
		}

		if string(b) == bb.String() {
			return nil
		}

		if !opts.YesToAll && !ask("Your webpack.config.js file is different from the latest Buffalo template.\nWould you like to replace yours with the latest template?") {
			fmt.Println("\tSkipping webpack.config.js")
			return nil
		}

		wf, err := os.Create("webpack.config.js")
		if err != nil {
			return err
		}
		_, err = wf.Write(bb.Bytes())
		if err != nil {
			return err
		}
		return wf.Close()
	}
}
