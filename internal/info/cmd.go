package info

import (
	"context"
	"os/exec"
	"time"

	"github.com/gobuffalo/clara/v2/genny/rx"
	"github.com/gobuffalo/cli/internal/genny/info"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/meta"
)

// Command instance to be used outside of this package.
var Command = newCommand()

type command struct {
	app   meta.App
	clara *rx.Options
	info  *info.Options
}

func (c command) Name() string {
	return "info"
}

func (c command) HelpText() string {
	return "Prints diagnostic information (useful for debugging)"
}

func (c command) Main(ctx context.Context, pwd string, args []string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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
		if err := run.WithNew(rx.New(c.clara)); err != nil {
			return err
		}
	}

	if err := run.WithNew(info.New(c.info)); err != nil {
		return err
	}

	return run.Run()
}

func newCommand() *command {
	app := meta.New(".")

	return &command{
		app: app,
		clara: &rx.Options{
			App: app,
		},
		info: &info.Options{
			App: app,
		},
	}
}
