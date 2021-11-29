package testhelpers

import (
	"os"
	"testing"
)

func RunWithinTempFolder(t *testing.T, fn func(t *testing.T)) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dir, err := os.MkdirTemp("", "buffalo-new-test-*")
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

	t.Cleanup(func() {
		err := os.Chdir(wd)
		if err != nil {
			t.Logf("error moving back to the original folder: %s", err)
		}
	})

}
