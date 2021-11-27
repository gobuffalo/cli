package testhelpers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"golang.org/x/mod/modfile"
)

// EnsureBuffaloCMD installs current version of buffalo, it should be run before
// calling RunBuffaloCMD to ensure this is the version of buffalo that is being tested.
func EnsureBuffaloCMD(t *testing.T) error {
	t.Helper()

	// Ensure we're on the
	ok, err := inCLISource()
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("not in the cli source folder")
	}

	ex := exec.Command("go", "install", "-tags", "sqlite", "github.com/gobuffalo/cli/cmd/buffalo")
	return ex.Run()
}

// Ensures that the current directory is the CLI source folder by
// checking its parent go.mod file says its github.com/gobuffalo/cli module.
func inCLISource() (bool, error) {
	wd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	mod := ""
	for {
		dat, err := ioutil.ReadFile(filepath.Join(wd, "go.mod"))
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

// RunBuffaloCMD is useful for integration tests where CMD would want
// to run a Buffalo command from the fully compiled binary.
func RunBuffaloCMD(t *testing.T, args []string) (string, error) {
	t.Helper()

	output := bytes.NewBufferString("")

	ex := exec.Command("buffalo")
	ex.Stdout = output
	ex.Stderr = output
	ex.Args = append(ex.Args, args...)
	err := ex.Run()

	return output.String(), err
}