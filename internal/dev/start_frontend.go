package dev

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/cli/internal/genny/assets/webpack"
	"github.com/gobuffalo/meta"
)

var StartFrontend = &startFrontend{}

type startFrontend struct {
	debug bool
}

func (ss startFrontend) Name() string {
	return "dev/start-frontend"
}

func (ss startFrontend) HelpText() string {
	return "Starts the frontent watcher."
}

func (ss *startFrontend) EnableDebug() {
	ss.debug = true
}

func (ss startFrontend) RunDevelopment(ctx context.Context, app meta.App, args []string) error {
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
			fmt.Println(err)
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
