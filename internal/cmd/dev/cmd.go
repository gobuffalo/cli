package dev

import (
	"strings"

	"github.com/gobuffalo/events"
	"github.com/spf13/cobra"
)

// debug flag to enable delve debugging
var debug bool = false

func Cmd() *cobra.Command {
	// Listen to events for event rewrite
	events.NamedListen("buffalo:dev", func(e events.Event) {
		if strings.HasPrefix(e.Kind, "refresh:") {
			e.Kind = strings.Replace(e.Kind, "refresh:", "buffalo:dev:", 1)
			events.Emit(e)
		}
	})

	// devCmd represents the dev command
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Run the Buffalo app in 'development' mode",
		Long: `Run the Buffalo app in 'development' mode.
This includes rebuilding the application when files change.
This behavior can be changed in .buffalo.dev.yml file.`,
		RunE: runE,
	}

	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "use delve to debug the app")

	return cmd
}
