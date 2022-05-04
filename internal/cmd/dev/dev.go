package dev

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/gobuffalo/cli/internal/genny/assets/webpack"
	rg "github.com/gobuffalo/cli/internal/genny/refresh"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/refresh/refresh"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func runE(c *cobra.Command, args []string) error {
	if runtime.GOOS == "windows" {
		color.NoColor = true
	}
	defer func() {
		msg := "There was a problem starting the dev server, Please review the troubleshooting docs: %s\n"
		cause := "Unknown"
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				cause = err.Error()
			}
		}
		logrus.Errorf(msg, cause)
	}()
	os.Setenv("GO_ENV", "development")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		return startDevServer(ctx, args)
	})

	wg.Go(func() error {
		app := meta.New(".")
		if !app.WithNodeJs {
			// No need to run dev script
			return nil
		}
		return runDevScript(ctx, app)
	})

	err := wg.Wait()
	if err != context.Canceled {
		return err
	}
	return nil
}

func runDevScript(ctx context.Context, app meta.App) error {
	tool := "yarnpkg"
	if !app.WithYarn {
		tool = "npm"
	}

	if _, err := exec.LookPath(tool); err != nil {
		return fmt.Errorf("could not find %s tool", tool)
	}

	// make sure that the node_modules folder is properly "installed"
	if _, err := os.Stat(filepath.Join(app.Root, "node_modules")); err != nil {
		cmd := exec.CommandContext(ctx, tool, "install")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	cmd := exec.CommandContext(ctx, tool, "run", "dev")
	if _, err := app.NodeScript("dev"); err != nil {
		// Fallback on legacy runner
		cmd = exec.CommandContext(ctx, webpack.BinPath, "--watch")
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return contextAwareRun(ctx, cmd.Run)
}

func startDevServer(ctx context.Context, args []string) error {
	app := meta.New(".")

	cfgFile := "./.buffalo.dev.yml"
	if _, err := os.Stat(cfgFile); err != nil {
		run := genny.WetRunner(ctx)
		err = run.WithNew(rg.New(&rg.Options{App: app}))
		if err != nil {
			return err
		}

		if err := run.Run(); err != nil {
			return err
		}
	}
	c := &refresh.Configuration{}
	if err := c.Load(cfgFile); err != nil {
		return err
	}
	c.Debug = debug

	bt := app.BuildTags("development")
	if len(bt) > 0 {
		c.BuildFlags = append(c.BuildFlags, "-tags", bt.String())
	}
	r := refresh.NewWithContext(c, ctx)
	r.CommandFlags = args

	return contextAwareRun(ctx, r.Start)
}

func contextAwareRun(ctx context.Context, f func() error) error {
	var out = make(chan error)

	go func() {
		out <- f()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-out:
		return err
	}
}
