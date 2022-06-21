package fix

import (
	"context"
	"flag"
	"io/ioutil"
	"os"

	"github.com/gobuffalo/cli/internal/genny/fix"
	"github.com/gobuffalo/cli/internal/runtime"
	"github.com/gobuffalo/genny/v2"
	"github.com/markbates/sigtx"
)

// Command instance of the fix command to be used by the CLI
var Command = &command{
	options: &fix.Options{},
}

type command struct {
	flagSet *flag.FlagSet
	options *fix.Options
}

func (c command) Name() string {
	return "fix"
}

func (c command) HelpText() string {
	return "Attempt to fix a Buffalo application's API to match version " + runtime.Version
}

func (c command) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet("fix", flag.ContinueOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(ioutil.Discard)
	}

	c.flagSet.BoolVar(&c.options.YesToAll, "y", false, "update all without asking for confirmation")

	return c.flagSet, nil
}

func (c command) Main(ctx context.Context, root string, args []string) error {
	ctx, cancel := sigtx.WithCancel(ctx, os.Interrupt)
	defer cancel()

	run := genny.WetRunner(ctx)
	if err := run.WithNew(fix.New(c.options)); err != nil {
		return err
	}
	return run.Run()
}
