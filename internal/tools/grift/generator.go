package grift

import (
	"context"

	"github.com/gobuffalo/cli/internal/genny/grift"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

// func init() {
// 	TaskCmd.Flags().BoolVarP(&taskOptions.dryRun, "dry-run", "d", false, "dry run of the generator")
// }

var Generator = &taskGenerator{}

type taskGenerator struct {
	dryRun bool
	*grift.Options
}

func (gg taskGenerator) Name() string {
	return "task"
}

func (gg taskGenerator) HelpText() string {
	return "Generate a grift task"
}

func (gg taskGenerator) Aliases() []string {
	return []string{"t", "grift"}
}

func (gg taskGenerator) Generate(ctx context.Context, pwd string, args []string) error {
	run := genny.WetRunner(context.Background())
	if gg.dryRun {
		run = genny.DryRunner(ctx)
	}

	opts := gg.Options
	opts.Args = args
	g, err := grift.New(opts)
	if err != nil {
		return err
	}
	run.With(g)

	g, err = gogen.Fmt(pwd)
	if err != nil {
		return err
	}
	run.With(g)

	return run.Run()
}
