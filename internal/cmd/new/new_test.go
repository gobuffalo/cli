package new

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	r := require.New(t)
	dir := os.TempDir()

	r.NoError(os.Chdir(dir))
	r.NoError(os.RemoveAll(filepath.Join(dir, "app")))

	tcases := []struct {
		name  string
		args  []string
		check func(*require.Assertions)
	}{
		{
			name: "skip docker",
			args: []string{"app", "--api", "--skip-docker", "-f", "--vcs", "none"},
			check: func(r *require.Assertions) {
				r.NoFileExists(filepath.Join("app", "Dockerfile"))
			},
		},

		{
			name: "docker there",
			args: []string{"app", "--api", "-f", "--vcs", "none"},
			check: func(r *require.Assertions) {
				r.FileExists(filepath.Join("app", "Dockerfile"))
			},
		},
	}

	for _, v := range tcases {
		t.Run(v.name, func(t *testing.T) {
			r := require.New(t)
			Cmd.SetArgs(v.args)

			r.NoError(Cmd.Execute())
			v.check(r)
		})
	}

}
