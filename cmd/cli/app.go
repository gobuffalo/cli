package cli

import (
	"context"
	"fmt"
	"sync"

	flag "github.com/spf13/pflag"
)

// app is the first point of entry for the CLI
// with the IO and plugins.
type App struct {
	*IO
	*flag.FlagSet

	plock   sync.RWMutex
	plugins Plugins
}

func (a *App) Run(ctx context.Context, pwd string, args []string) error {
	a.plock.RLock()
	defer a.plock.RUnlock()

	if len(args) == 0 {
		// TODO: Print help
		return nil
	}

	// Find the command that should run.
	cmd := a.plugins.FindCommand(args[0])
	if cmd == nil {
		// TODO: Print help
		return fmt.Errorf("no command found for %s", args[0])
	}

	// Set the IO if the command supports it
	if is, ok := cmd.(IOSetter); ok {
		is.SetIO(a.Stdin(), a.Stdout(), a.Stderr())
	}

	// Call the flag parsing on the command and
	// update the args after the flag parsing.
	if fp, ok := cmd.(FlagParser); ok {
		var err error
		args, err = fp.ParseFlags(args)
		if err != nil {
			return err
		}
	}

	return cmd.Run(ctx, pwd, args)
}
