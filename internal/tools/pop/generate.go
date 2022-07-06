package pop

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/gobuffalo/cli/cmd/cli/help"
	"github.com/gobuffalo/cli/internal/tools/pop/generators"
)

var Generate = &generate{
	generators: Generators{
		generators.Config,
		generators.Fizz,
		generators.SQL,
		generators.Model,
	},
}

type generate struct {
	generators Generators
}

func (g generate) Name() string {
	return "generate"
}

func (g generate) Usage() string {
	return "buffalo pop generate <generator> [flags] [options]"
}

func (g generate) Aliases() []string {
	return []string{"g"}
}

func (g generate) HelpText() string {
	return "Generates config, model, and migrations files."
}

func (g *generate) LongHelpText() string {
	output := bytes.NewBuffer([]byte{})
	w := tabwriter.NewWriter(output, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "Generators:\n")
	for _, v := range g.generators {
		fmt.Fprintf(w, "  %s\t%s\n", v.Name(), v.HelpText())
	}

	w.Flush()
	return output.String()
}

func (g generate) Help(ctx context.Context, args []string) error {
	cc := g.generators.Find(args[0])
	hh, ok := cc.(help.Helper)
	if !ok || len(args) == 1 {
		return help.Specific(os.Stdout, cc)
	}

	// If the command implements Helper
	// the command itself will take care
	// of printing the help with the args.
	return hh.Help(ctx, args[1:])
}

func (g generate) PopMain(ctx context.Context, pwd string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify the generator name")
	}

	gg := g.generators.Find(args[0])
	if gg == nil {
		return fmt.Errorf("did not find generator %s", args[0])
	}

	return gg.PopGenerate(ctx, pwd, args[1:])
}
