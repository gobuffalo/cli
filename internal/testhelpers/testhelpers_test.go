//go:build integration
// +build integration

package testhelpers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/cli/internal/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestEnsureBuffaloCMD(t *testing.T) {
	r := require.New(t)

	path := filepath.Join(os.TempDir(), "buffalointegrationtests")
	os.Remove(path)

	t.Log("Original >>", path)

	r.NoError(testhelpers.EnsureBuffaloCMD(t))
	r.FileExists(path)
}
