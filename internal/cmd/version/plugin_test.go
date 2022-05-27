package version_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gobuffalo/cli/internal/cmd/version"
	"github.com/gobuffalo/cli/internal/runtime"
)

func TestCommand(t *testing.T) {
	out := bytes.NewBuffer([]byte{})

	// Using the SetIO method to set the stdout and stderr
	version.Plugin.SetIO(nil, out, out)

	err := version.Plugin.Run(nil, "", []string{"version"})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), "Buffalo CLI version is: ") {
		t.Errorf("expected %s to contain %s, it did not", out.String(), "Buffalo CLI version is: ")
	}

	if !strings.Contains(out.String(), runtime.Version) {
		t.Errorf("expected %s to contain %s, it did not", out.String(), runtime.Version)
	}
}

func TestCommandJSON(t *testing.T) {
	out := bytes.NewBuffer([]byte{})

	// Using the SetIO method to set the stdout and stderr
	version.Plugin.SetIO(nil, out, out)
	args, err := version.Plugin.ParseFlags([]string{"version", "--json"})
	err = version.Plugin.Run(nil, "", args)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), `"version": "`) {
		t.Errorf("expected %s to contain %s, it did not", out.String(), `"version": "`)
	}

	if !strings.Contains(out.String(), runtime.Version) {
		t.Errorf("expected %s to contain %s, it did not", out.String(), runtime.Version)
	}
}
