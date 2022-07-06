package generate

import (
	"context"
	"fmt"
	"io"

	flag "github.com/spf13/pflag"

	"github.com/gobuffalo/cli/internal/genny/actions"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/logger"
)

var ActionGenerator = &actionGenerator{
	options: &actions.Options{},
}

type actionGenerator struct {
	flagSet *flag.FlagSet
	options *actions.Options

	dryRun  bool
	verbose bool
}

func (ag actionGenerator) Usage() string {
	return "generate action [name] [handler name...]"
}

func (ag actionGenerator) Name() string {
	return "action"
}

func (ag actionGenerator) HelpText() string {
	return "Generate new action(s)"
}

func (ag actionGenerator) Aliases() []string {
	return []string{"a"}
}

func (ag *actionGenerator) ParseFlags(args []string) (*flag.FlagSet, error) {
	if ag.flagSet == nil {
		ag.flagSet = flag.NewFlagSet("action", flag.ContinueOnError)
		ag.flagSet.Usage = func() {}
		ag.flagSet.SetOutput(io.Discard)
	}

	ag.flagSet.BoolVar(&ag.dryRun, "dry-run", false, "Runs the generator without writing any files.")
	ag.flagSet.BoolVar(&ag.verbose, "verbose", false, "Prints more verbose output.")

	_ = ag.flagSet.Parse(args)

	return ag.flagSet, nil
}

func (ag actionGenerator) Generate(ctx context.Context, pwd string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a name")
	}

	ag.options.Name = args[0]
	if len(args) == 1 {
		return fmt.Errorf("you must provide at least one action name")
	}

	ag.options.Actions = args[1:]
	run := genny.WetRunner(ctx)

	if ag.dryRun {
		run = genny.DryRunner(ctx)
	}

	if ag.verbose {
		run.Logger = logger.New(logger.DebugLevel)
	}

	if err := run.WithNew(actions.New(ag.options)); err != nil {
		return err
	}

	return run.Run()
}
