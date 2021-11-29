// +build integration

package testhelpers_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gobuffalo/cli/internal/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestEnsureBuffaloCMD(t *testing.T) {
	r := require.New(t)

	binary := "buffalointegrationtests"
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}

	path := filepath.Join(os.TempDir(), binary)
	os.Remove(path)

	r.NoError(testhelpers.EnsureBuffaloCMD(t))
	r.FileExists(path)
}
