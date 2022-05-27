package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/gobuffalo/cli/cmd/cli"
	"github.com/gobuffalo/cli/internal/cmd/version"
)

func TestVersionPlugin(t *testing.T) {
	out := bytes.NewBuffer([]byte{})
	app := cli.New(
		cli.WithIO(&cli.IO{
			Out: out,
			Err: out,
		}),
		cli.WithPlugins(version.Plugin),
	)

	err := app.Run(context.Background(), "", []string{"version"})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), "Buffalo CLI version is: ") {
		t.Errorf("expected %s to contain %s, it did not", out.String(), "Buffalo CLI version is: ")
	}
}

func TestNoCommandFound(t *testing.T) {
	out := bytes.NewBuffer([]byte{})
	app := cli.New(
		cli.WithIO(&cli.IO{
			Out: out,
			Err: out,
		}),
		cli.WithPlugins(version.Plugin),
	)

	err := app.Run(context.Background(), "", []string{"nocommand"})
	if err == nil {
		t.Fatal(err)
	}
}
