package help_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	flag "github.com/spf13/pflag"

	"github.com/gobuffalo/cli/cmd/cli/help"
	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

type testCommand string

func (t testCommand) Aliases() []string {
	return []string{
		"tc",
		"test-command",
	}
}

func (t testCommand) Name() string {
	return string(t)
}

func (t testCommand) HelpText() string {
	return fmt.Sprintf("runs the %v thing basically", t)
}

func (t testCommand) LongHelpText() string {
	return fmt.Sprintf("Long text for the command that runs the %v thing. We could list here subcommands and steps.", t)
}

func (t testCommand) ParseFlags(args []string) (*flag.FlagSet, error) {
	var value string
	var toggle bool

	fs := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
	fs.StringVar(&value, "output", "plain", "Output format (string)")
	fs.BoolVar(&toggle, "toggle", false, "A toggle flag (boolean)")

	// This is to keep it silent
	fs.SetOutput(bytes.NewBuffer([]byte{}))
	fs.Usage = func() {}

	// Ignore the error we don't care if any error happens while parsing.
	_ = fs.Parse(args)

	return fs, nil
}

func (t testCommand) Main(ctx context.Context, pwd string, args []string) error {
	return nil
}

func TestHelpCommand(t *testing.T) {
	hc := help.Command{
		Commands: plugin.Commands{
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

		if !bytes.Contains(out.Bytes(), []byte("Usage: buffalo [command] [options]")) {
			t.Fatalf("expected output to contain 'Usage: buffalo [command] [options]'")
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

		if !bytes.Contains(out.Bytes(), []byte("Usage: buffalo [command] [options]")) {
			t.Fatalf("expected output to contain 'Usage: buffalo [command] [options]'")
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
			"Usage: buffalo test [options]",
			"--output",
			"--toggle",
		}

		for _, v := range contents {
			if bytes.Contains(out.Bytes(), []byte(v)) {
				continue
			}

			t.Fatalf("expected output to contain '%s'", v)
		}

		if !bytes.Contains(out.Bytes(), []byte("Long text for the command that runs")) {
			t.Fatalf("expected output to contain 'Long text for the command that runs'")
		}

	})

	t.Run("specific command and subcommands", func(t *testing.T) {
		out := bytes.NewBuffer([]byte{})
		hc.IO.Out = out

		err := hc.Main(context.Background(), "", []string{"test", "something"})
		if err != nil {
			t.Fatalf("error running help: %v", err)
		}

		if !bytes.Contains(out.Bytes(), []byte("Usage: buffalo test [options]")) {
			t.Fatalf("expected output to contain 'Usage: buffalo test [options]'")
		}
	})

}

type simplePlugin string

func (tc simplePlugin) Name() string {
	return string(tc)
}

func TestReceivePlugins(t *testing.T) {
	plugins := plugin.Plugins{
		testCommand("test"),
		simplePlugin("simple"),
	}

	hc := &help.Command{}
	hc.Receive(plugins)

	if len(hc.Commands) != 1 {
		t.Fatalf("expected 1 command, got %v", len(hc.Commands))
	}
}
