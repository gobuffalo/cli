package generate

import (
	"context"
	"flag"
	"fmt"

	"github.com/gobuffalo/cli/internal/genny/actions"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/logger"
)

var ActionGenerator = &actionGenerator{
	flagSet: flag.NewFlagSet("action", flag.ContinueOnError),
}

type actionGenerator struct {
	flagSet *flag.FlagSet
	options *actions.Options

	dryRun  bool
	verbose bool
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

// ActionCmd is the cmd that generates actions.
// var ActionCmd = &cobra.Command{
// 	Use:     "action [name] [handler name...]",
// 	Aliases: []string{"a", "actions"},
// 	Short:   "Generate new action(s)",
// 	RunE: func(cmd *cobra.Command, args []string) error {

// 	},
// }
