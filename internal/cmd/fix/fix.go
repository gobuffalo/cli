package fix

import (
	"context"
	"os"

	"github.com/gobuffalo/cli/internal/genny/fix"
	"github.com/gobuffalo/genny/v2"
	"github.com/markbates/sigtx"
	"github.com/spf13/cobra"
)

// run all compatible checks
func RunE(cmd *cobra.Command, args []string) error {
	ctx, cancel := sigtx.WithCancel(context.Background(), os.Interrupt)
	defer cancel()

	opts := &fix.Options{
		YesToAll: yesToAll,
	}

	run := genny.WetRunner(ctx)
	if err := run.WithNew(fix.New(opts)); err != nil {
		return err
	}
	return run.Run()
}
