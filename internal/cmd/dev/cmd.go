package dev

import (
	"strings"

	"github.com/gobuffalo/events"
	"github.com/spf13/cobra"
)

var (
	debug bool
)

func Cmd() *cobra.Command {
	// Listen to events for event rewrite
	events.NamedListen("buffalo:dev", func(e events.Event) {
		if strings.HasPrefix(e.Kind, "refresh:") {
			e.Kind = strings.Replace(e.Kind, "refresh:", "buffalo:dev:", 1)
			events.Emit(e)
		}
	})

	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "use delve to debug the app")

	return cmd
}
