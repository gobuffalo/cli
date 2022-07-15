package cli

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/cli/cmd/cli/help"
)

// DefaultApp is an instance of the CLI application
// loaded with `default` plugins. The `NewApp` function
// could be used to create a custom instance of the CLI
// with custom plugins.
var DefaultApp = &App{
	help:    help.Command,
	plugins: append(basePlugins, defaultPlugins...),

	overriders: []overrider{
		projectOverrider,
	},
}

// projectOverrider returns an exec.Cmd if the user is overriding the CLI within the
// same module. It looks within the current directory for the cmd/buffalo/main.go
// and attempts to run that file.
func projectOverrider(pwd string, args []string) (*exec.Cmd, string) {
	if _, err := os.Stat(filepath.Join(pwd, "cmd", "buffalo", "main.go")); err != nil {
		return nil, ""
	}

	cmd := exec.Command("go")
	cmd.Args = append(cmd.Args, "run", filepath.Join(pwd, "cmd", "buffalo", "main.go"))
	cmd.Args = append(cmd.Args, args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdout

	return cmd, "cmd/buffalo"
}

// Here we take care of looking for CLI overriders
// Overriders are go files to run instead of the CLI,
// The two main use cases are:
//  1. Running codebase specific CLI (PWD/cmd/buffalo/main.go) ✅
//  2. Running user specific CLI ($HOME/buffalo/cmd/main.go)
