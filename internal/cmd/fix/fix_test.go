// +build integration

package fix_test

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/cli/internal/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestFix(t *testing.T) {
	r := require.New(t)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))

	t.Cleanup(func() {
		buffaloBin := filepath.Join(build.Default.GOPATH, "bin", "buffalo")
		if err := os.Remove(buffaloBin); err != nil {
			t.Logf("failed to delete buffalo binary: %s", buffaloBin)
		}
	})

	versions := []string{
		"v0.18.0",
		"v0.17.7",
		"v0.16.27",
		"v0.15.5",
	}

	tt := []struct {
		newargs []string
		appname string
	}{
		{
			newargs: []string{"new", "api", "-f", "--api", "--vcs", "none"},
			appname: "api",
		},
		{
			newargs: []string{"new", "web", "-f", "--vcs", "none"},
			appname: "web",
		},
	}

	for _, version := range versions {
		r.NoError(testhelpers.InstallOldBuffaloCMD(t, version))

		for _, tc := range tt {
			testname := fmt.Sprintf("%s - %s", tc.appname, version)

			t.Run(testname, func(t *testing.T) {
				testhelpers.RunWithinTempFolder(t, func(t *testing.T) {
					r := require.New(t)

					out, err := exec.Command("buffalo", tc.newargs...).CombinedOutput()
					t.Log(string(out))
					r.NoError(err)

					r.NoError(os.Chdir(tc.appname))

					output, err := testhelpers.RunBuffaloCMD(t, []string{"fix", "-y"})
					t.Log(output)
					r.NoError(err)

					output, err = testhelpers.RunBuffaloCMD(t, []string{"build"})
					t.Log(output)
					r.NoError(err)
				})
			})
		}
	}
}
