package generate

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/tabwriter"

	flag "github.com/spf13/pflag"

	"github.com/gobuffalo/cli/cmd/cli/clio"
	"github.com/gobuffalo/cli/cmd/cli/help"
	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

// Command instance to be used outside of the
// generate package
var Command = &generate{
	IO:         &clio.IO{},
	generators: Generators{},
}

// Generate command is the root command for the
// different generators that CLI could have.
type generate struct {
	*clio.IO

	generators Generators
}

func (g generate) Name() string {
	return "generate"
}

func (g generate) Aliases() []string {
	return []string{"g", "gen"}
}

func (g generate) HelpText() string {
	return "Generates code blocks in the source code, aiming to accelerate software development."
}

func (g generate) Usage() string {
	return "buffalo generate [generator] [options]"
}

func (c generate) ValidateWorkDir(wd string) (bool, error) {
	return plugin.ValidateBuffaloRoot(wd)
}

func (g *generate) LongHelpText() string {
	if len(g.generators) == 0 {
		return "No generators registered.\n"
	}

	buf := bytes.NewBuffer([]byte{})
	w := tabwriter.NewWriter(buf, 0, 0, 3, ' ', 0)

	w.Write([]byte("Registered Generators\n"))
	for _, gg := range g.generators {
		aliases := []string{}
		if gh, ok := gg.(plugin.Aliaser); ok {
			aliases = gh.Aliases()
		}

		fmt.Fprintf(w, "%s\t%+v\t%s\n", gg.Name(), strings.Join(aliases, ", "), gg.HelpText())
	}

	w.Flush()

	return buf.String()
}

func (g *generate) ParseFlags(args []string) (*flag.FlagSet, error) {
	if len(args) > 0 {
		args = args[1:]
	}

	for _, v := range g.generators {
		if pf, ok := v.(clio.FlagParser); ok {
			pf.ParseFlags(args)
		}
	}

	return flag.NewFlagSet("generate", flag.ContinueOnError), nil
}

func (g *generate) Help(ctx context.Context, args []string) error {
	if len(args) < 1 {
		fmt.Fprintln(g.Stdout(), g.LongHelpText())

		return nil
	}

	// Find the generator
	// Print its help text
	gg := g.generators.Find(args[0])
	if gg == nil {
		fmt.Fprintf(g.Stdout(), "Error: No generator found for '%v'\n\n", args[0])
		fmt.Fprintln(g.Stdout(), g.LongHelpText())

		return nil
	}

	usage := "buffalo generate " + gg.Name()
	if hh, ok := gg.(help.Usager); ok {
		usage = hh.Usage()
	}

	fmt.Fprintf(g.Stdout(), "Usage: %v\n\n", usage)
	fmt.Fprintln(g.Stdout(), gg.HelpText())
	if ht, ok := gg.(help.LongHelpTexter); ok {
		fmt.Fprintln(g.Stdout(), ht.LongHelpText())
	}

	if fl, ok := gg.(clio.FlagParser); ok {
		fl, _ := fl.ParseFlags([]string{})

		fmt.Fprint(g.Stdout(), "\nFlags:\n")
		fl.VisitAll(func(ff *flag.Flag) {
			fmt.Fprintf(g.Stdout(), "--%v\t%v\n", ff.Name, ff.Usage)
		})
	}

	return nil
}

// Receive the plugins and select the ones that implement the
// Generator interface, these will be stored within the
// generators variable.
func (g *generate) Receive(plugins plugin.Plugins) {
	g.generators = Generators{}

	for _, p := range plugins {
		gg, ok := p.(Generator)
		if !ok {
			continue
		}

		g.generators = append(g.generators, gg)
	}
}

func (g generate) Main(ctx context.Context, pwd string, args []string) error {
	if len(args) < 1 {
		return g.Help(ctx, args)
	}

	gg := g.generators.Find(args[0])
	if gg == nil {
		return g.Help(ctx, args)
	}

	return gg.Generate(ctx, pwd, args)
}
