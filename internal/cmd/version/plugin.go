package version

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gobuffalo/cli/internal/runtime"

	flag "github.com/spf13/pflag"
)

// Plugin for the version
var Plugin = &command{
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

func (c *command) Run(ctx context.Context, pwd string, args []string) error {
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

func (c *command) ParseFlags(args []string) ([]string, error) {
	c.flags.BoolVar(&c.json, "json", false, "output in JSON format")
	c.flags.Parse(args)

	return c.flags.Args(), nil
}

func (c *command) SetIO(stdin io.Reader, stdout, stderr io.Writer) {
	c.stdout = stdout
	c.stderr = stderr
}
