package testhelpers

import (
	"os"
	"testing"
)

// RunWithinTempFolder runs the given function on a temporary folder
// and returns to the original working directory afterwards.
func RunWithinTempFolder(t *testing.T, fn func(t *testing.T)) {
	t.Helper()

	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err := os.Chdir(original)
		if err != nil {
			t.Logf("error moving back to the original folder: %s", err)
		}
	})

	dir, err := os.MkdirTemp("", "buffalo-integration-test-*")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("failed to delete temporary directory: %s", dir)
		}
	})

	err = os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}

	fn(t)
}
