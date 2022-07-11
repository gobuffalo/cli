package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/gobuffalo/cli/cmd/cli"
)

type plx string

func (p plx) Name() string {
	return string(p)
}

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

func TestAdd(t *testing.T) {
	app := cli.DefaultApp
	bb := &bytes.Buffer{}
	app.IO.Out = bb

	app.Add(
		plx("fake/plugin"),
		plx("fake/plugin-2"),
	)

	err := app.Main(context.TODO(), "", []string{"plugins"})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(bb.String())

	for _, v := range []string{
		"Loaded default CLI plugins",
		"Plugins loaded (4)",
		"fake/plugin",
		"fake/plugin-2",
	} {
		if !strings.Contains(bb.String(), v) {
			t.Fatalf("expected to contain '%v'", v)
		}
	}

}
