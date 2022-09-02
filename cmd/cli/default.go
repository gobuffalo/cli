package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/cli/cmd/cli/help"
	"github.com/gobuffalo/cli/cmd/cli/plugin"
	"github.com/gobuffalo/cli/internal/build"
	"github.com/gobuffalo/cli/internal/destroy"
	"github.com/gobuffalo/cli/internal/dev"
	"github.com/gobuffalo/cli/internal/fix"
	"github.com/gobuffalo/cli/internal/generate"
	"github.com/gobuffalo/cli/internal/info"
	"github.com/gobuffalo/cli/internal/new"
	"github.com/gobuffalo/cli/internal/routes"
	"github.com/gobuffalo/cli/internal/setup"
	"github.com/gobuffalo/cli/internal/test"
	"github.com/gobuffalo/cli/internal/tools/bzr"
	"github.com/gobuffalo/cli/internal/tools/frontend"
	"github.com/gobuffalo/cli/internal/tools/git"
	"github.com/gobuffalo/cli/internal/tools/grift"
	"github.com/gobuffalo/cli/internal/tools/pop"
	"github.com/gobuffalo/cli/internal/version"
)

var (
	basePlugins = plugin.Plugins{
		help.Command,
		plugin.List,
	}

	// defaultPlugins used by the CLI. This could be used as base set
	// when customizing the CLI.
	defaultPlugins = plugin.Plugins{
		// Top level commands
		test.Command,
		version.Command,
		grift.Command,
		routes.Command,
		setup.Command, // TODO: DOCS: Document how to wire setuppers here.
		info.Command,
		generate.Command, // TODO: DOCS: Document how to add generators
		fix.Command,
		new.Command,
		destroy.Command,
		build.Command,
		dev.Command, // TODO: DOCS: Document how to add dev plugins
		pop.Command, // TODO: This needs to live outside of the CLI package and into the pop/buffalo-pop package.

		pop.Create,
		pop.Drop,
		pop.Fix,
		pop.Migrate,
		pop.Reset,
		pop.Generate,
		pop.Schema,

		// Generators
		generate.ActionGenerator,
		generate.MailerGenerator,
		generate.ResourceGenerator,
		grift.Generator,

		//Destroyers
		destroy.ActionDestroyer,
		destroy.MailerDestroyer,
		destroy.ResourceDestroyer,

		// Setup plugins
		grift.SetupSeedDatabase,
		pop.Setup,
		frontend.Setup,
		setup.Test,

		// Development plugins
		dev.SetupDevelopment,
		dev.StartServer,
		dev.StartFrontend,

		git.VersionRunner,
		bzr.VersionRunner,
	}
)

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

// userOverrider returns an exec.Cmd if the user is overriding the CLI
// on $HOME/.buffalo/cmd/main.go
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
