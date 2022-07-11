package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/gobuffalo/cli/cmd/cli"
)

func TestClear(t *testing.T) {
	app := cli.DefaultApp
	bb := &bytes.Buffer{}
	app.IO.Out = bb

	app.Clear()

	err := app.Main(context.TODO(), "", []string{"plugins"})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(bb.String())

	for _, v := range []string{
		"Loaded default CLI plugins",
		"Plugins loaded (2)",
	} {
		if !strings.Contains(bb.String(), v) {
			t.Fatalf("expected to contain '%v'", v)
		}
	}

}
