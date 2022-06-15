package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/gobuffalo/cli/cmd/cli"
)

func TestHelpCommand(t *testing.T) {

	hc := cli.HelpCommand{
		Commands: cli.Commands{
			testCommand("test"),
			testCommand("other"),
		},
	}

	t.Run("plain help invoked", func(t *testing.T) {
		out := bytes.NewBuffer([]byte{})
		hc.IO.Out = out

		err := hc.Main(context.Background(), "", []string{})
		if err != nil {
			t.Fatalf("error running help: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("Usage: depbot [command] [options]")) {
			t.Fatalf("expected output to contain 'Usage: depbot [command] [options]'")
		}

		for _, v := range hc.Commands {
			if !bytes.Contains(out.Bytes(), []byte(fmt.Sprintf("%v", v.Name()))) {
				t.Fatalf("expected output to contain '%v'", v.Name())
			}
		}
	})

	t.Run("unexisting command invoked on help", func(t *testing.T) {
		out := bytes.NewBuffer([]byte{})
		hc.IO.Out = out

		err := hc.Main(context.Background(), "", []string{"unexisting"})
		if err != nil {
			t.Fatalf("error running help: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("Error: did not find `unexisting` command")) {
			t.Fatalf("Expected to print \"Error: did not find `unexisting` command\"")
		}

		if !bytes.Contains(out.Bytes(), []byte("Usage: depbot [command] [options]")) {
			t.Fatalf("expected output to contain 'Usage: depbot [command] [options]'")
		}
	})

	t.Run("specific command", func(t *testing.T) {
		out := bytes.NewBuffer([]byte{})
		hc.IO.Out = out

		err := hc.Main(context.Background(), "", []string{"test"})
		if err != nil {
			t.Fatalf("error running help: %v", err)
		}

		contents := []string{
			"Usage: depbot test [options]",
			"--output",
			"--toggle",
		}

		for _, v := range contents {
			if bytes.Contains(out.Bytes(), []byte(v)) {
				continue
			}

			t.Fatalf("expected output to contain '%s'", v)
		}
	})

	t.Run("specific command and subcommands", func(t *testing.T) {
		out := bytes.NewBuffer([]byte{})
		hc.IO.Out = out

		err := hc.Main(context.Background(), "", []string{"test", "something"})
		if err != nil {
			t.Fatalf("error running help: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("Usage: depbot test [options]")) {
			t.Fatalf("expected output to contain 'Usage: depbot test [options]'")
		}
	})

}
