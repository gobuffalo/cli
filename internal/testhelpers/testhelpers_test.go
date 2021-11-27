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

	r.NoError(testhelpers.EnsureBuffaloCMD(t))

	homeDir, err := os.UserHomeDir()
	r.NoError(err)

	r.FileExists(filepath.Join(homeDir, "buffalointegrationtests"))
}
