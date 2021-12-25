// +build integration

package fix_test

import (
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/cli/internal/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestFix_v0_18_0(t *testing.T) {
	r := require.New(t)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))
	r.NoError(testhelpers.InstallOldBuffaloCMD(t, "v0.18.0"))

	t.Cleanup(func() {
		buffaloBin := filepath.Join(build.Default.GOPATH, "bin", "buffalo")
		if err := os.Remove(buffaloBin); err != nil {
			t.Logf("failed to delete buffalo binary: %s", buffaloBin)
		}
	})

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

	for _, tc := range tt {
		t.Run(tc.appname, func(t *testing.T) {
			testhelpers.RunWithinTempFolder(t, func(t *testing.T) {
				r := require.New(t)

				ex := exec.Command("buffalo", tc.newargs...)
				ex.Stdout = os.Stdout
				ex.Stderr = os.Stderr
				r.NoError(ex.Run())

				r.NoError(os.Chdir(tc.appname))

				out, err := testhelpers.RunBuffaloCMD(t, []string{"fix", "-y"})
				t.Log(out)
				r.NoError(err)

				out, err = testhelpers.RunBuffaloCMD(t, []string{"build"})
				t.Log(out)
				r.NoError(err)
			})
		})
	}
}

func TestFix_v0_17_7(t *testing.T) {
	r := require.New(t)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))
	r.NoError(testhelpers.InstallOldBuffaloCMD(t, "v0.17.7"))

	t.Cleanup(func() {
		buffaloBin := filepath.Join(build.Default.GOPATH, "bin", "buffalo")
		if err := os.Remove(buffaloBin); err != nil {
			t.Logf("failed to delete buffalo binary: %s", buffaloBin)
		}
	})

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

	for _, tc := range tt {
		t.Run(tc.appname, func(t *testing.T) {
			testhelpers.RunWithinTempFolder(t, func(t *testing.T) {
				r := require.New(t)

				ex := exec.Command("buffalo", tc.newargs...)
				ex.Stdout = os.Stdout
				ex.Stderr = os.Stderr
				r.NoError(ex.Run())

				r.NoError(os.Chdir(tc.appname))

				out, err := testhelpers.RunBuffaloCMD(t, []string{"fix", "-y"})
				t.Log(out)
				r.NoError(err)

				out, err = testhelpers.RunBuffaloCMD(t, []string{"build"})
				t.Log(out)
				r.NoError(err)
			})
		})
	}
}

func TestFix_v0_16_27(t *testing.T) {
	r := require.New(t)
	r.NoError(testhelpers.EnsureBuffaloCMD(t))
	r.NoError(testhelpers.InstallOldBuffaloCMD(t, "v0.16.27"))

	t.Cleanup(func() {
		buffaloBin := filepath.Join(build.Default.GOPATH, "bin", "buffalo")
		if err := os.Remove(buffaloBin); err != nil {
			t.Logf("failed to delete buffalo binary: %s", buffaloBin)
		}
	})

	tt := []struct {
		newargs      []string
		appname      string
		replacements map[string]string // replace existing files that are not fixed by buffalo fix but need updates to compile
	}{
		{
			newargs: []string{"new", "api", "-f", "--api", "--vcs", "none"},
			appname: "api",
		},
		{
			newargs: []string{"new", "web", "-f", "--vcs", "none"},
			appname: "web",
			replacements: map[string]string{
				"assets/js/application.js": `
require("expose-loader?exposes=$,jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");
require("@fortawesome/fontawesome-free/js/all.js");
require("jquery-ujs/src/rails.js");

$(() => {

});`,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.appname, func(t *testing.T) {
			testhelpers.RunWithinTempFolder(t, func(t *testing.T) {
				r := require.New(t)

				ex := exec.Command("buffalo", tc.newargs...)
				ex.Stdout = os.Stdout
				ex.Stderr = os.Stderr
				r.NoError(ex.Run())

				r.NoError(os.Chdir(tc.appname))

				for k, v := range tc.replacements {
					r.NoError(os.WriteFile(k, []byte(v), 0644))
				}

				out, err := testhelpers.RunBuffaloCMD(t, []string{"fix", "-y"})
				t.Log(out)
				r.NoError(err)

				out, err = testhelpers.RunBuffaloCMD(t, []string{"build"})
				t.Log(out)
				r.NoError(err)
			})
		})
	}
}
