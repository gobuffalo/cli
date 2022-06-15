package cli

import (
	"github.com/gobuffalo/cli/internal/routes"
	"github.com/gobuffalo/cli/internal/test"
	"github.com/gobuffalo/cli/internal/tools/grift"
	"github.com/gobuffalo/cli/internal/version"
)

var (
	// DefaultApp is an instance of the CLI application
	// loaded with `default` plugins. The `NewApp` function
	// could be used to create a custom instance of the CLI
	// with custom plugins.
	DefaultApp = NewApp(
		test.Command,

		// The version command
		version.Command,

		// The task command
		grift.Command,

		// The routes command
		routes.Command,
	)
)
