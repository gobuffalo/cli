package pop

import (
	"context"

	"github.com/gobuffalo/cli/cmd/cli/plugin"
	"github.com/gobuffalo/pop/v6"
)

// Subcommand is the interface of the plugins that could be hooked into
// the pop command.
type Subcommand interface {
	plugin.Plugin

	Run(context.Context, []pop.ConnectionDetails, []string) error
}
