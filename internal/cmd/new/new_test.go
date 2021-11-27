package new

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/cli/internal/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	r := require.New(t)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))

	dir, err := os.MkdirTemp("", "buffalo-new-test-*")
	r.NoError(err)
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("failed to delete temporary directory: %s", dir)
		}
	})

	r.NoError(os.Chdir(dir))

	tcases := []struct {
		name  string
		args  []string
		check func(*require.Assertions, string, error)
	}{
		{
			name: "no application name",
			args: []string{"new"},
			check: func(r *require.Assertions, out string, err error) {
				r.Error(err)
				r.Contains(out, "you must enter a name for your new application")
			},
		},
		{
			name: "skip docker",
			args: []string{"new", "nodocker", "--api", "--skip-docker", "-f", "--vcs", "none"},
			check: func(r *require.Assertions, out string, err error) {
				r.NoError(err)
				r.NoFileExists(filepath.Join("nodocker", "Dockerfile"))
			},
		},

		{
			name: "docker there",
			args: []string{"new", "wdocker", "--api", "-f", "--vcs", "none"},
			check: func(r *require.Assertions, out string, err error) {
				r.NoError(err)
				r.FileExists(filepath.Join("wdocker", "Dockerfile"))
			},
		},

		{
			name: "invalid db type",
			args: []string{"new", "api", "--api", "-f", "--db-type", "a"},
			check: func(r *require.Assertions, out string, err error) {
				r.Error(err)
				r.Contains(out, `unknown dialect`)
			},
		},

		{
			name: "forbidden application name",
			args: []string{"new", "buffalo", "-f", "--api"},
			check: func(r *require.Assertions, out string, err error) {
				r.Error(err)
				r.Contains(out, `name buffalo is not allowed, try a different application name`)
			},
		},
	}

	for _, v := range tcases {
		t.Run(v.name, func(t *testing.T) {
			r := require.New(t)
			out, err := testhelpers.RunBuffaloCMD(t, v.args)
			v.check(r, out, err)
		})
	}
}
