package dev

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/fatih/color"
	"github.com/gobuffalo/cli/cmd/cli/plugin"
	"github.com/gobuffalo/events"
	"github.com/gobuffalo/meta"
	"golang.org/x/sync/errgroup"
)

var Command = &command{}

type command struct {
	flagSet *flag.FlagSet

	debug bool

	setuppers []DevelopmentSetupper
	runners   []DevelopmentRunner
}

func (c *command) Name() string {
	return "dev"
}

func (c *command) Aliases() []string {
	return []string{"dev", "develop", "serve", "server"}
}

func (c *command) HelpText() string {
	return "Run the Buffalo app in 'development' mode"
}

func (c *command) LongHelpText() string {
	return `This includes rebuilding the application when files change.
This behavior can be changed in .buffalo.dev.yml file.`
}

func (c *command) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(io.Discard)
	}

	c.flagSet.BoolVar(&c.debug, "debug", false, "use delve to debug the app")

	_ = c.flagSet.Parse(args)

	if c.debug {
		for _, v := range c.runners {
			v.EnableDebug()
		}
	}

	return c.flagSet, nil
}

func (c *command) Receive(pls plugin.Plugins) {
	for _, v := range pls {
		if s, ok := v.(DevelopmentSetupper); ok {
			c.setuppers = append(c.setuppers, s)
		}

		if r, ok := v.(DevelopmentRunner); ok {
			c.runners = append(c.runners, r)
		}
	}
}

func (c *command) Main(ctx context.Context, pwd string, args []string) error {
	defer func() {
		cause := "Unknown"
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				cause = err.Error()
			}
		}

		fmt.Printf("There was a problem starting the dev server, Please review the troubleshooting docs: %s\n", cause)
	}()

	// Listen to events for event rewrite
	events.NamedListen("buffalo:dev", func(e events.Event) {
		if !strings.HasPrefix(e.Kind, "refresh:") {
			return
		}

		e.Kind = strings.Replace(e.Kind, "refresh:", "buffalo:dev:", 1)
		events.Emit(e)
	})

	color.NoColor = (runtime.GOOS == "windows")

	// Running setupers
	for _, s := range c.setuppers {
		err := s.SetupDevelopment(ctx, pwd, args)
		if err == nil {
			continue
		}

		return err
	}

	app := meta.New(pwd)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg, ctx := errgroup.WithContext(ctx)
	for _, r := range c.runners {
		wg.Go(func(r DevelopmentRunner) func() error {
			return func() error {
				return r.RunDevelopment(ctx, app, args)
			}
		}(r))
	}

	err := wg.Wait()
	if err != context.Canceled {
		return err
	}

	return nil
}
