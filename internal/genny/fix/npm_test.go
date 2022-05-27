package fix

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_AddPackageJSONScripts_AddScriptSection(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		App: meta.App{
			Root:        ".",
			WithWebpack: true,
		},
	}

	run := gentest.NewRunner()
	r.NoError(run.File(genny.NewFileS("package.json", "{}")))

	run.WithRun(AddPackageJSONScripts(opts))
	r.NoError(run.Run())

	results := run.Results()
	f := results.Files[0]
	r.Equal("package.json", f.Name())

	packageJSON := map[string]map[string]string{}
	r.NoError(json.NewDecoder(f).Decode(&packageJSON))
	r.EqualValues(defaultScripts, packageJSON["scripts"])
}

func Test_AddPackageJSONScripts_AddMissingScripts(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		App: meta.App{
			Root:        ".",
			WithWebpack: true,
		},
	}

	tt := []struct {
		Name    string
		Scripts map[string]string
	}{
		{
			Name:    "EmptyScripts",
			Scripts: map[string]string{},
		},
		{
			Name: "MissingDev",
			Scripts: map[string]string{
				"build": "echo 'hello'",
				"start": "echo 'hello'",
			},
		},
		{
			Name: "MissingBuild",
			Scripts: map[string]string{
				"dev":   "echo 'hello'",
				"start": "echo 'hello'",
			},
		},
		{
			Name: "NoMissingScripts",
			Scripts: map[string]string{
				"dev":   "echo 'hello'",
				"build": "echo 'hello'",
				"start": "echo 'hello'",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			run := gentest.NewRunner()

			packageJSON := map[string]map[string]string{
				"scripts": tc.Scripts,
			}
			f := genny.NewFile("package.json", nil)
			r.NoError(json.NewEncoder(f).Encode(packageJSON))
			r.NoError(run.File(f))

			run.WithRun(AddPackageJSONScripts(opts))
			r.NoError(run.Run())

			results := run.Results()
			f = results.Files[0]
			r.Equal("package.json", f.Name())

			r.NoError(json.NewDecoder(f).Decode(&packageJSON))

			for k, v := range packageJSON["scripts"] {
				if _, ok := tc.Scripts[k]; !ok {
					r.Equalf(defaultScripts[k], v, "missing default script: %s", k)
				} else {
					r.Equalf(tc.Scripts[k], v, "script %s has been altered", k)
				}
			}
		})
	}
}

func Test_PackageJSONCheck_NoOverwriteExisting(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		App: meta.App{
			Root:        ".",
			WithWebpack: true,
		},
	}

	run := gentest.NewRunner()
	bb, err := defaultPackageJson(opts.App)
	r.NoError(err)
	fileContents := bb.String()
	r.NoError(run.File(genny.NewFileS("package.json", fileContents)))

	run.WithRun(PackageJSONCheck(opts))
	r.NoError(run.Run())

	results := run.Results()
	f := results.Files[0]
	r.Equal("package.json", f.Name())
	r.Equal(fileContents, f.String())
}

func Test_PackageJSONCheck_UpdatingFileAndClearNodeModules(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		App: meta.App{
			Root:        ".",
			WithWebpack: true,
		},
		YesToAll: true,
	}

	bb, err := defaultPackageJson(opts.App)
	r.NoError(err)
	fileContents := bb.String()

	tt := []struct {
		Name     string
		WithYarn bool
		Command  string
	}{
		{
			Name:     "WithYarn",
			WithYarn: true,
			Command:  "yarn install",
		},
		{
			Name:     "WithoutYarn",
			WithYarn: false,
			Command:  "npm install",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			run := gentest.NewRunner()
			r.NoError(run.File(genny.NewFileS("package.json", "{}")))
			r.NoError(run.File(genny.NewFileS("node_modules/test.txt", "this file should be deleted")))

			opts.App.WithYarn = tc.WithYarn
			run.WithRun(PackageJSONCheck(opts))
			r.NoError(run.Run())

			results := run.Results()

			r.Len(results.Files, 1, "node_modules has not been deleted")

			f := results.Files[0]
			r.Equal("package.json", f.Name())
			r.Equal(fileContents, f.String())

			r.Len(results.Commands, 1, "command has not been run")
			c := results.Commands[0]
			r.Equal(tc.Command, strings.Join(c.Args, " "))
		})
	}
}
