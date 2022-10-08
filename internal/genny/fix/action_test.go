package fix

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_ReplaceAppOnce(t *testing.T) {
	r := require.New(t)

	tt := []struct {
		Name string
		OK   bool
	}{
		{"buffalo0_11", false},
		{"buffaloPre0_18api", true},
		{"buffaloPre0_18web", true},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			run := gentest.NewRunner()

			err := run.Disk.AddFS(os.DirFS(filepath.Join("_fixtures", tc.Name)))
			r.NoError(err)

			opts := &Options{
				App: meta.Named("coke", "."),
			}
			g := ReplaceAppOnce(opts)
			run.WithRun(g)

			if tc.OK {
				r.NoError(run.Run())
				f, err := run.FindFile("actions/app.go")
				r.NoError(err)

				//fmt.Println("XXX =====", f.String())
				r.Contains(f.String(), "appOnce", "files in vendor directory should not be changed")
			} else {
				r.Error(run.Run())
			}
		})
	}
}
