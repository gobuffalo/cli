package testhelpers

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/mod/modfile"
)

var buffaloVersions = []struct {
	installUrl string
	apiUrl     string
}{
	{
		"github.com/gobuffalo/cli/cmd/buffalo",
		"https://api.github.com/repos/gobuffalo/cli/releases/tags/%s",
	},
	{
		"github.com/gobuffalo/buffalo/buffalo",
		"https://api.github.com/repos/gobuffalo/buffalo/releases/tags/%s",
	},
}

func getBuffaloUrl(version string) (string, error) {
	for _, v := range buffaloVersions {
		url := fmt.Sprintf(v.apiUrl, version)
		resp, err := http.Get(url)
		if err != nil {
			return "", fmt.Errorf("failed to fetch url %s: %w", v, err)
		}
		if resp.StatusCode == http.StatusOK {
			return fmt.Sprintf("%s@%s", v.installUrl, version), nil
		}
	}

	return "", fmt.Errorf("unknown gobuffalo cli version %s", version)
}

// InstallOldBuffaloCMD installs a specific version of buffalo for
// integration tests.
func InstallOldBuffaloCMD(t *testing.T, version string) error {
	t.Helper()

	url, err := getBuffaloUrl(version)
	if err != nil {
		return err
	}
	t.Log("installing", url)

	ex := exec.Command("go")
	ex.Args = append(ex.Args,
		"install",
		"-tags",
		"sqlite",
		url,
	)

	ex.Stdout = os.Stdout
	ex.Stderr = os.Stderr
	return ex.Run()
}

// EnsureBuffaloCMD installs current version of buffalo, it should be run before
// calling RunBuffaloCMD to ensure this is the version of buffalo that is being tested.
func EnsureBuffaloCMD(t *testing.T) error {
	t.Helper()

	ok, err := inCLISource()
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("not in the cli source folder")
	}

	ex := exec.Command("go")
	ex.Args = append(
		ex.Args,
		"build",
		"-tags",
		"sqlite",
		"-o",
		testingBinaryLocation(t),
		"github.com/gobuffalo/cli/cmd/buffalo",
	)

	ex.Stdout = os.Stdout
	ex.Stderr = os.Stderr
	return ex.Run()
}

// RunBuffaloCMD is useful for integration tests where CMD would want
// to run a Buffalo command from the fully compiled binary.
func RunBuffaloCMD(t *testing.T, args []string) (string, error) {
	t.Helper()

	ex := exec.Command(testingBinaryLocation(t))
	ex.Args = append(ex.Args, args...)
	output, err := ex.CombinedOutput()

	return string(output), err
}

// testingBinaryLocation returns the location of the testing binary which is
// set to be the user home folder on a file named `buffalointegrationtests`.
func testingBinaryLocation(t *testing.T) string {
	t.Helper()

	binary := "buffalointegrationtests"
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}

	return filepath.Join(os.TempDir(), binary)
}

// inCLISource ensures that the current directory is the CLI source folder by
// checking its parent go.mod file says its github.com/gobuffalo/cli module.
func inCLISource() (bool, error) {
	wd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	mod := ""
	for {
		dat, err := os.ReadFile(filepath.Join(wd, "go.mod"))
		if err != nil {
			wd = filepath.Dir(wd)
			if wd == "/" {
				break
			}

			continue
		}

		f, err := modfile.Parse("go.mod", dat, nil)
		if err != nil {
			return false, err
		}

		mod = f.Module.Mod.Path
		break
	}

	result := mod == "github.com/gobuffalo/cli"
	return result, nil
}
