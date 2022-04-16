package fix

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func TestMoveMain(t *testing.T) {
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

			g := MoveMain(opts)
			run.WithRun(g)

			r.NoError(run.Run())

			results := run.Results()
			_, err = results.Find("cmd/app/main.go")
			r.NoError(err)

			_, err = results.Find("main.go")
			r.Error(err)
		})
	}
}
