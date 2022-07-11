package plugin

import (
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/gobuffalo/cli/cmd/cli/clio"
)

var List = &listCommand{}

type listCommand struct {
	clio.IO

	plugins Plugins
}

func (lc listCommand) Name() string {
	return "plugins"
}

func (lc listCommand) HelpText() string {
	return "List plugins loaded to the CLI."
}

func (lc *listCommand) Receive(pls Plugins) {
	lc.plugins = pls
}

func (lc listCommand) Main(ctx context.Context, pwd string, args []string) error {
	type helpTexter interface{ HelpText() string }

	descFn := func(v Plugin, t string) string {
		var helpText string
		if v, ok := v.(helpTexter); ok {
			helpText = v.HelpText()
		}

		return fmt.Sprintf("  %s\t%s\t%s\n", v.Name(), t, helpText)
	}

	fmt.Fprintln(lc.Stdout(), "Loaded default CLI plugins (buffalo binary).")
	fmt.Fprintf(lc.Stdout(), "\nPlugins loaded (%v):\n", len(lc.plugins))

	w := tabwriter.NewWriter(lc.Stdout(), 0, 0, 3, ' ', 0)
	for _, v := range CommandsFrom(lc.plugins) {
		fmt.Fprintf(w, descFn(v, "[command]"))
	}

	fmt.Fprintf(w, "  \t\t\n")

	for _, v := range lc.plugins {
		if _, ok := v.(Command); ok {
			continue
		}

		fmt.Fprintf(w, descFn(v, "[other]"))
	}

	w.Flush()

	return nil
}
