package fix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/cli/internal/genny/assets/webpack"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/meta"
)

// default scripts for package.json
var defaultScripts = map[string]string{
	"dev":   "webpack --watch",
	"build": "webpack -p --progress",
}

// AddPackageJSONScripts rewrites the package.json file
// to add dev and build scripts if there are missing.
func AddPackageJSONScripts(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		if !opts.App.WithWebpack {
			return nil
		}

		fmt.Println("~~~ Patching package.json to add dev and build scripts ~~~")
		f, err := r.FindFile("package.json")
		if err != nil {
			return err
		}

		needRewrite := false
		packageJSON := map[string]interface{}{}
		if err := json.NewDecoder(f).Decode(&packageJSON); err != nil {
			return fmt.Errorf("could not rewrite package.json: %s", err.Error())
		}

		if _, ok := packageJSON["scripts"]; !ok {
			needRewrite = true
			packageJSON["scripts"] = defaultScripts
		} else {
			scripts, ok := packageJSON["scripts"].(map[string]interface{})
			if !ok {
				return fmt.Errorf("could not rewrite package.json: invalid scripts section")
			}

			// Add missing scripts
			for k, v := range defaultScripts {
				if _, ok := scripts[k]; !ok {
					scripts[k] = v
					needRewrite = true
				}
			}
			packageJSON["scripts"] = scripts
		}

		if !needRewrite {
			fmt.Println("~~~ package.json doesn't need to be patched, skipping. ~~~")
			return nil
		}

		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(packageJSON); err != nil {
			return fmt.Errorf("could not rewrite package.json: %w", err)
		}

		return nil
	}
}

// PackageJSONCheck will compare the current default Buffalo
// package.json against the applications package.json. If they are
// different you have the option to overwrite the existing package.json
// file with the new one.
func PackageJSONCheck(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		if !opts.App.WithWebpack {
			return nil
		}

		fmt.Println("~~~ Checking package.json ~~~")
		bb, err := defaultPackageJson(opts.App)
		if err != nil {
			return err
		}

		f, err := r.FindFile("package.json")
		if err != nil {
			return err
		}

		if f.String() == bb.String() {
			return nil
		}

		if !opts.YesToAll && !ask("Your package.json file is different from the latest Buffalo template.\nWould you like to REPLACE yours with the latest template?") {
			fmt.Println("\tskipping package.json")
			return nil
		}

		_, err = f.Write(bb.Bytes())
		if err != nil {
			return err
		}

		base := "node_modules"
		for _, f := range r.Disk.Files() {
			rel, err := filepath.Rel(base, f.Name())
			if err != nil {
				return err
			}

			if strings.HasPrefix(rel, "..") {
				continue
			}

			if err := r.Disk.Delete(f.Name()); err != nil {
				return err
			}
		}

		if opts.App.WithYarn {
			return r.Exec(exec.Command("yarnpkg", "install"))
		}

		return r.Exec(exec.Command("npm", "install"))
	}
}

func defaultPackageJson(app meta.App) (*bytes.Buffer, error) {
	templates, err := webpack.Templates()
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("package.json").ParseFS(templates, "package.json.tmpl")
	if err != nil {
		return nil, err
	}

	bb := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(bb, "package.json.tmpl", map[string]interface{}{
		"opts": &webpack.Options{
			App: app,
		},
	})
	return bb, err
}
