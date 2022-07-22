package plugin

import (
	"context"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/gobuffalo/cli/cmd/cli/clio"
)

// List command instance to be used on the app.
var List = &listCommand{}

// listCommand lists the plugins received by the CLI
// The CLI will pass those plugins to this command via
// the PluginReceiver interface which is implemented
// on this command.
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

func (lc *listCommand) SetIO(stdout io.Writer, stderr io.Writer, stdin io.Reader) {
	lc.IO.Out = stdout
	lc.IO.Err = stderr
	lc.IO.In = stdin
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

	fmt.Fprintf(lc.Stdout(), "Plugins loaded (%v):\n", len(lc.plugins))

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