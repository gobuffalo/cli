package fix

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_Depreciations_mainGo(t *testing.T) {
	r := require.New(t)

	tt := []struct {
		Name     string
		warnings []string
	}{
		{
			Name:     "buffalo0_11",
			warnings: []string{"app.Start has been removed in v0.11.0. Use app.Serve Instead. [main.go]"},
		},
		{
			Name:     "buffaloPre0_18api",
			warnings: []string{},
		},
		{
			Name:     "buffaloPre0_18web",
			warnings: []string{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			run := gentest.NewRunner()

			err := run.Disk.AddFS(os.DirFS(filepath.Join("_fixtures", tc.Name)))
			r.NoError(err)

			opts := &Options{
				App: meta.Named("coke", "."),
			}
			g := DeprecationsCheck(opts)
			run.WithRun(g)

			r.NoError(run.Run())

			results := run.Results()
			_, err = results.Find("main.go")
			r.NoError(err)

			r.ElementsMatch(opts.warnings, tc.warnings)
		})
	}
}

func Test_Depreciations_ReplacePackr(t *testing.T) {
	r := require.New(t)

	tt := []struct {
		Name        string
		contains    map[string][]string
		notContains map[string][]string
	}{
		{
			Name: "buffaloPre0_18api",
			contains: map[string][]string{
				"actions/app.go": {
					"coke/locales",
					"i18n.New(locales.FS(), \"en-US\")",
				},
			},
			notContains: map[string][]string{
				"actions/app.go": {
					"packr.New",
				},
				"actions/render.go": {
					"packr.New",
				},
			},
		},
		{
			Name: "buffaloPre0_18web",
			contains: map[string][]string{
				"actions/app.go": {
					"coke/locales",
					"coke/public",
					"net/http",
					"app.ServeFiles(\"/\", http.FS(public.FS()))",
					"i18n.New(locales.FS(), \"en-US\")",
				},
				"actions/render.go": {
					"coke/public",
					"coke/templates",
					"AssetsFS: public.FS()",
					"TemplatesFS: templates.FS()",
				},
			},
			notContains: map[string][]string{
				"actions/render.go": {
					"packr.New",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			run := gentest.NewRunner()

			err := run.Disk.AddFS(os.DirFS(filepath.Join("_fixtures", tc.Name)))
			r.NoError(err)

			opts := &Options{
				App: meta.Named("coke", "."),
			}
			g := DeprecationsCheck(opts)
			run.WithRun(g)

			r.NoError(run.Run())
			results := run.Results()

			clean := func(s string) string {
				s = strings.TrimSpace(s)
				s = strings.ReplaceAll(s, "\n", "")
				s = strings.ReplaceAll(s, "\t", "")
				s = strings.ReplaceAll(s, "\r", "")

				spaces := regexp.MustCompile(`\s+`)
				return spaces.ReplaceAllString(s, " ")
			}

			for file, contains := range tc.contains {
				f, err := results.Find(file)
				r.NoError(err)

				for _, c := range contains {
					r.Contains(clean(f.String()), clean(c))
				}
			}

			for file, notContains := range tc.notContains {
				f, err := results.Find(file)
				r.NoError(err)

				for _, c := range notContains {
					r.NotContains(clean(f.String()), clean(c))
				}
			}
		})
	}
}
