package fix

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_Imports(t *testing.T) {
	r := require.New(t)

	tt := []struct {
		Name string
	}{
		{
			Name: "buffalo0_11",
		},
		{
			Name: "buffaloPre0_18api",
		},
		{
			Name: "buffaloPre0_18web",
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
			g := ReplaceOldImports(opts)
			run.WithRun(g)

			r.NoError(run.Run())
			results := run.Results()

			for _, f := range results.Files {
				if filepath.Ext(f.Name()) != ".go" {
					continue
				}

				if f.Name() == "vendor/models_test.go" {
					r.Contains(f.String(), strconv.Quote("github.com/gobuffalo/suite"), "files in vendor directory should not be changed")
					continue
				}

				for k := range replace {
					r.NotContainsf(f.String(), strconv.Quote(k), "%s should not have %s", f.Name(), k)
				}
			}
		})
	}
}
