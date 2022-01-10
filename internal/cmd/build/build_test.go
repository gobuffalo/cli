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

	tt := []struct {
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

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			testhelpers.RunWithinTempFolder(t, func(t *testing.T) {
				r := require.New(t)
				out, err := testhelpers.RunBuffaloCMD(t, tc.newargs)
				t.Log(out)
				r.NoError(err)

				r.NoError(os.Chdir(tc.appname))

				// NOTE: I think adding this to here is totally fine since
				// this is "integration" test. However, the original reason
				// I added it now is to prevent build failure when there is
				// no subdir under templates directory (go:embed * */*)
				out, err = testhelpers.RunBuffaloCMD(t, tc.resourceargs)
				t.Log(out)
				r.NoError(err)

				out, err = testhelpers.RunBuffaloCMD(t, []string{"build"})
				t.Log(out)
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
	testhelpers.RunWithinTempFolder(t, func(t *testing.T) {
		out, err := testhelpers.RunBuffaloCMD(t, []string{"new", "noassets", "-f", "--skip-webpack", "--vcs", "none"})
		t.Log(out)
		r.NoError(err)

		t.Cleanup(func() {
			os.RemoveAll("noassets")
		})

		r.NoError(os.Chdir("noassets"))

		out, err = testhelpers.RunBuffaloCMD(t, []string{"g", "resource", "phone", "model"})
		t.Log(out)
		r.NoError(err)

		out, err = testhelpers.RunBuffaloCMD(t, []string{"build", "--extract-assets"})
		t.Log(out)
		r.NoError(err)

		r.FileExists(filepath.Join("bin", "assets.zip"))
	})
}
