package cli

import (
	"github.com/gobuffalo/cli/internal/routes"
	"github.com/gobuffalo/cli/internal/setup"
	"github.com/gobuffalo/cli/internal/test"
	"github.com/gobuffalo/cli/internal/tools/frontend"
	"github.com/gobuffalo/cli/internal/tools/grift"
	"github.com/gobuffalo/cli/internal/tools/pop"
	"github.com/gobuffalo/cli/internal/version"
)

var (
	// DefaultApp is an instance of the CLI application
	// loaded with `default` plugins. The `NewApp` function
	// could be used to create a custom instance of the CLI
	// with custom plugins.
	DefaultApp = NewApp(
		// Top level commands
		test.Command,
		version.Command,
		grift.Command,
		routes.Command,
		setup.Command,

		// Setup plugins
		grift.SetupSeedDatabase,
		pop.Setup,
		frontend.Setup,
		setup.Test,
	)
)
