package cli

import (
	"github.com/gobuffalo/cli/internal/cmd/version"
)

var (
	// App is the default CLI app with the default plugin set.
	DefaultApp = New(
		WithPlugins(version.Plugin),
	)
)
