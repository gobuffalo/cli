package setup

import "github.com/gobuffalo/meta"

// Setuper is the type of those plugins that are used
// to setup the application.
type Setuper interface {
	// Name of the Setuper, useful for help texts and
	// error messages.
	Name() string

	// Method that will be called to setup the application.
	// receives meta.App with the application details.
	Setup(meta.App) error
}
