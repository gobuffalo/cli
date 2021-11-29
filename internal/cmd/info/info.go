package info

import (
	"context"
	"os/exec"
	"time"

	"github.com/gobuffalo/clara/v2/genny/rx"
	"github.com/gobuffalo/cli/internal/genny/info"
	"github.com/gobuffalo/genny/v2"
	"github.com/spf13/cobra"
)

func runE(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	run := genny.WetRunner(ctx)

	_, err := run.LookPath("clara")
	if err == nil {
		// use the clara binary if available
		run.WithRun(func(r *genny.Runner) error {
			return r.Exec(exec.Command("clara"))
		})
	} else {
		// no clara binary, so use the one bundled with buffalo
		copts := infoOptions.Clara
		if err := run.WithNew(rx.New(copts)); err != nil {
			return err
		}
	}

	iopts := infoOptions.Info
	if err := run.WithNew(info.New(iopts)); err != nil {
		return err
	}

	return run.Run()
}
