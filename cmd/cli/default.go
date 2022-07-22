package cli

import (
	"fmt"
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
		userOverrider,
	},
}

// projectOverrider returns an exec.Cmd if the user is overriding the CLI within the
// same module. It looks within the current directory for the cmd/buffalo/main.go
// and attempts to run that file.
func projectOverrider(pwd string) (*exec.Cmd, string) {
	if _, err := os.Stat(filepath.Join(pwd, "cmd", "buffalo", "main.go")); err != nil {
		return nil, ""
	}

	args := os.Args
	if len(args) > 0 {
		args = os.Args[1:]
	}

	cmd := exec.Command("go")
	cmd.Args = append(cmd.Args, "run", filepath.Join(pwd, "cmd", "buffalo", "main.go"))
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin

	return cmd, "cmd/buffalo"
}

func userOverrider(pwd string) (*exec.Cmd, string) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, ""
	}

	path := filepath.Join(home, ".buffalo", "cmd", "main.go")
	if _, err = os.Stat(path); err != nil {
		return nil, ""
	}

	args := os.Args
	if len(args) > 0 {
		args = os.Args[1:]
	}

	cmd := exec.Command("go")
	cmd.Args = append(cmd.Args, "run", path)
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = os.Stdin

	return cmd, fmt.Sprintf("%v/buffalo", home)
}
