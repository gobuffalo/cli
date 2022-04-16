package fix

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func TestRefresh(t *testing.T) {
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

			g := Refresh(opts)
			run.WithRun(g)

			r.NoError(run.Run())
			results := run.Results()

			f, err := results.Find(".buffalo.dev.yml")

			r.NoError(err)
			r.Contains(f.String(), "./cmd/app")
		})
	}
}
