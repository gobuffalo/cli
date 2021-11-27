package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/cli/internal/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	r := require.New(t)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))

	tcases := []struct {
		name    string
		newargs []string
		appname string
	}{
		{
			name:    "nominal",
			newargs: []string{"new", "nominal", "-f", "--skip-webpack", "--vcs", "none"},
			appname: "nominal",
		},
		{
			name:    "api",
			newargs: []string{"new", "api", "-f", "--api", "--vcs", "none"},
			appname: "api",
		},
		{
			name:    "sqlite",
			newargs: []string{"new", "sqlite", "-f", "--skip-webpack", "--db-type=sqlite3", "--vcs", "none"},
			appname: "sqlite",
		},
	}

	dir, err := os.MkdirTemp("", "buffalo-build-test-*")
	r.NoError(err)
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("failed to delete temporary directory: %s", dir)
		}
	})

	for _, v := range tcases {
		t.Run(v.name, func(tx *testing.T) {
			r := require.New(tx)

			r.NoError(os.Chdir(dir))

			out, err := testhelpers.RunBuffaloCMD(t, v.newargs)
			tx.Log(out)
			r.NoError(err)

			r.NoError(os.Chdir(filepath.Join(dir, v.appname)))

			out, err = testhelpers.RunBuffaloCMD(t, []string{"build"})
			tx.Log(out)
			r.NoError(err)
		})
	}
}
