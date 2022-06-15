package version

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/gobuffalo/cli/internal/runtime"
)

// Command for the version
var Command = &command{
	flags: flag.NewFlagSet("version", flag.ContinueOnError),
}

type command struct {
	flags *flag.FlagSet
	json  bool

	// io fields
	stdout io.Writer
	stderr io.Writer
}

func (c *command) Name() string {
	return "version"
}

func (c *command) ParseFlags(args []string) (*flag.FlagSet, error) {
	c.flags.BoolVar(&c.json, "json", false, "Print information in json format")
	c.flags.Parse(args)

	return c.flags, nil
}

func (c *command) Main(ctx context.Context, pwd string, args []string) error {
	if c.json {
		enc := json.NewEncoder(c.stdout)
		enc.SetIndent("", "    ")

		return enc.Encode(runtime.BuildInfo{
			Version: runtime.Version,
		})
	}

	fmt.Fprintf(c.stdout, "Buffalo CLI version is: %s\n", runtime.Version)
	return nil
}

func (c *command) SetIO(stdin io.Reader, stdout, stderr io.Writer) {
	c.stdout = stdout
	c.stderr = stderr
}

func (c command) HelpText() string {
	return "Prints the version of the CLI in plan and JSON formats."
}
