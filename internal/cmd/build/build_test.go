// +build integration

package build_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gobuffalo/cli/internal/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	r := require.New(t)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))

	tcases := []struct {
		name         string
		newargs      []string
		resourceargs []string
		appname      string
	}{
		{
			name:         "nominal",
			newargs:      []string{"new", "nominal", "-f", "--skip-webpack", "--vcs", "none"},
			resourceargs: []string{"g", "resource", "phone", "model"},
			appname:      "nominal",
		},
		{
			name:         "api",
			newargs:      []string{"new", "api", "-f", "--api", "--vcs", "none"},
			resourceargs: []string{"g", "resource", "phone", "model"},
			appname:      "api",
		},
		{
			name:         "sqlite",
			newargs:      []string{"new", "sqlite", "-f", "--skip-webpack", "--db-type=sqlite3", "--vcs", "none"},
			resourceargs: []string{"g", "resource", "phone", "model"},
			appname:      "sqlite",
		},
		{
			name:         "skipop",
			newargs:      []string{"new", "skipop", "-f", "--skip-pop", "--vcs", "none"},
			resourceargs: []string{"g", "action", "phone", "new"},
			appname:      "skipop",
		},
	}

	for _, v := range tcases {
		t.Run(v.name, func(tx *testing.T) {
			testhelpers.RunWithinTempFolder(tx, func(tt *testing.T) {
				r := require.New(tt)
				out, err := testhelpers.RunBuffaloCMD(tt, v.newargs)
				tt.Log(out)
				r.NoError(err)

				os.Chdir(v.appname)

				// NOTE: I think adding this to here is totally fine since
				// this is "integration" test. However, the original reason
				// I added it now is to prevent build failure when there is
				// no subdir under templates directory (go:embed * */*)
				out, err = testhelpers.RunBuffaloCMD(t, v.resourceargs)
				tx.Log(out)
				r.NoError(err)

				out, err = testhelpers.RunBuffaloCMD(tt, []string{"build"})
				tt.Log(out)
				r.NoError(err)
			})
		})
	}
}

func TestBuildNoAssets(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping this test on windows temporarily")
	}

	r := require.New(t)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))

	testhelpers.RunWithinTempFolder(t, func(tt *testing.T) {
		out, err := testhelpers.RunBuffaloCMD(tt, []string{"new", "noassets", "-f", "--skip-webpack", "--vcs", "none"})
		tt.Log(out)
		r.NoError(err)

		tt.Cleanup(func() {
			os.RemoveAll("noassets")
		})

		os.Chdir("noassets")

		out, err = testhelpers.RunBuffaloCMD(t, []string{"g", "resource", "phone", "model"})
		tt.Log(out)
		r.NoError(err)

		out, err = testhelpers.RunBuffaloCMD(tt, []string{"build", "--extract-assets"})
		tt.Log(out)
		r.NoError(err)

		r.FileExists(filepath.Join("bin", "assets.zip"))
	})
}
