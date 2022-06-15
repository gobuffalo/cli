package plugin_test

import (
	"bytes"
	"context"
	"flag"
	"fmt"
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
