//go:build integration
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

	for _, v := range tcases {
		t.Run(v.name, func(tx *testing.T) {
			r := require.New(tx)
			wd, err := os.Getwd()
			r.NoError(err)
			defer os.Chdir(wd)

			dir := os.TempDir()
			os.Chdir(filepath.Join(dir))

			out, err := testhelpers.RunBuffaloCMD(tx, v.newargs)
			tx.Log(out)
			r.NoError(err)

			os.Chdir(filepath.Join(dir, v.appname))
			out, err = testhelpers.RunBuffaloCMD(tx, []string{"build"})
			tx.Log(out)
			r.NoError(err)
		})
	}
}

func TestBuildNoAssets(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping this test on windows temporarily")
	}

	r := require.New(t)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))

	wd, err := os.Getwd()
	r.NoError(err)
	defer os.Chdir(wd)

	dir := os.TempDir()
	os.Chdir(dir)

	out, err := testhelpers.RunBuffaloCMD(t, []string{"new", "noassets", "-f", "--skip-webpack", "--vcs", "none"})
	t.Log(out)
	r.NoError(err)

	os.Chdir(filepath.Join(dir, "noassets"))
	out, err = testhelpers.RunBuffaloCMD(t, []string{"build", "--extract-assets"})
	t.Log(out)
	r.NoError(err)

	r.FileExists(filepath.Join("bin", "assets.zip"))
}
