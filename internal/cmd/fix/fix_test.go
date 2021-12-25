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

	tt := []struct {
		version      string
		newargs      []string
		appname      string
		replacements map[string]string // replace existing files that are not fixed by buffalo fix but need updates to compile
	}{
		{
			version: "v0.18.0",
			newargs: []string{"new", "api", "-f", "--api", "--vcs", "none"},
			appname: "api",
		},
		{
			version: "v0.18.0",
			newargs: []string{"new", "web", "-f", "--vcs", "none"},
			appname: "web",
		},

		{
			version: "v0.17.7",
			newargs: []string{"new", "api", "-f", "--api", "--vcs", "none"},
			appname: "api",
		},
		{
			version: "v0.17.7",
			newargs: []string{"new", "web", "-f", "--vcs", "none"},
			appname: "web",
		},

		{
			version: "v0.16.27",
			newargs: []string{"new", "api", "-f", "--api", "--vcs", "none"},
			appname: "api",
		},
		{
			version: "v0.16.27",
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
		t.Run(fmt.Sprintf("%s %s", tc.appname, tc.version), func(t *testing.T) {
			testhelpers.RunWithinTempFolder(t, func(t *testing.T) {
				r := require.New(t)
				r.NoError(testhelpers.InstallOldBuffaloCMD(t, tc.version))

				ex := exec.Command("buffalo", tc.newargs...)
				ex.Stdout = os.Stdout
				ex.Stderr = os.Stderr
				ex.Run()

				r.NoError(os.Chdir(tc.appname))

				for k, v := range tc.replacements {
					r.NoError(os.WriteFile(k, []byte(v), 0644))
				}

				r.NoError(testhelpers.RunBuffaloCMD(t, []string{"fix", "-y"}))
				r.NoError(testhelpers.RunBuffaloCMD(t, []string{"build"}))
			})
		})
	}
}
