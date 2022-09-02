package clio_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/gobuffalo/cli/cmd/cli/clio"
)

func TestIO(t *testing.T) {
	io := clio.IO{}

	if io.Stdout() != os.Stdout {
		t.Errorf("expected Stdout to default to os.Stdout")
	}

	if io.Stderr() != os.Stderr {
		t.Errorf("expected Stderr to default to os.Stderr")
	}

	if io.Stdin() != os.Stdin {
		t.Errorf("expected Stdin to default to os.Stdin")
	}

	out := bytes.NewBuffer([]byte{})
	in := strings.NewReader("")

	io.Err = out
	io.In = in
	io.Out = out

	if io.Stderr() != out {
		t.Errorf("expected Err to be set to out")
	}

	if io.Stdout() != out {
		t.Errorf("expected Out to be set to out")
	}

	if io.Stdin() != in {
		t.Errorf("expected In to be set to in")
	}
}
