package fix

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/cli/internal/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestFix(t *testing.T) {
	r := require.New(t)
	pwd, err := os.Getwd()
	r.NoError(err)
	r.NoError(testhelpers.InstallBuffaloCMD(t, "v0.17.5"))

	dir, err := os.MkdirTemp("", "buffalo-fix-test-*")
	r.NoError(err)
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("failed to delete temporary directory: %s", dir)
		}
	})

	tt := []struct {
		name    string
		newargs []string
		appname string
	}{
		{
			name:    "api",
			newargs: []string{"new", "api", "-f", "--api", "--vcs", "none"},
			appname: "api",
		},
		{
			name:    "web",
			newargs: []string{"new", "web", "-f", "--skip-webpack", "--vcs", "none"},
			appname: "web",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(tx *testing.T) {
			r := require.New(tx)

			r.NoError(os.Chdir(dir))

			out, err := testhelpers.RunBuffaloCMD(t, tc.newargs)
			tx.Log(out)
			r.NoError(err)
		})
	}

	r.NoError(os.Chdir(pwd))
	r.NoError(testhelpers.EnsureBuffaloCMD(t))

	for _, tc := range tt {
		t.Run(tc.name, func(tx *testing.T) {
			r := require.New(tx)

			r.NoError(os.Chdir(filepath.Join(dir, tc.appname)))

			out, err := testhelpers.RunBuffaloCMD(t, []string{"fix", "-y"})
			tx.Log(out)
			r.NoError(err)

			out, err = testhelpers.RunBuffaloCMD(t, []string{"build"})
			tx.Log(out)
			r.NoError(err)
		})
	}

	// TODO: is a successful build after fix enough for the test or should we check that the now fixed application actually matches a newly generated one?
}
