package generate

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/gobuffalo/cli/cmd/cli/help"
	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

// Command instance to be used outside of the
// generate package
var Command = &generate{}

// Generate command is the root command for the
// different generators that CLI could have.
type generate struct {
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

func (g *generate) LongHelpText() string {
	buf := bytes.NewBuffer([]byte{})
	w := tabwriter.NewWriter(buf, 0, 0, 3, ' ', 0)

	w.Write([]byte("Registered Generators\n"))
	for _, gg := range g.generators {
		aliases := []string{}
		if gh, ok := gg.(plugin.Aliaser); ok {
			aliases = gh.Aliases()
		}

		line := fmt.Sprintf("%s\t%+v\t%s\n", gg.Name(), strings.Join(aliases, ", "), gg.HelpText())
		w.Write([]byte(line))
	}

	w.Flush()

	return buf.String()
}

func (g *generate) Help(ctx context.Context, args []string) error {
	// Find the generator
	// Print its help text
	gg := g.generators.Find(args[0])
	if gg == nil {
		fmt.Printf("Error: No generator found for '%v'\n\n", args[0])
		fmt.Println(g.LongHelpText())

		return nil
	}

	usage := "buffalo generate " + gg.Name()
	if hh, ok := gg.(help.Usager); ok {
		usage = hh.Usage()
	}

	fmt.Printf("Usage: %v\n\n", usage)
	fmt.Println(gg.HelpText())
	if ht, ok := gg.(help.LongHelpTexter); ok {
		fmt.Println(ht.LongHelpText())
	}

	return nil
}

func (g *generate) Receive(plugins plugin.Plugins) {
	for _, p := range plugins {
		gg, ok := p.(Generator)
		if !ok {
			continue
		}

		g.generators = append(g.generators, gg)
	}
}

func (g generate) Main(ctx context.Context, pwd string, args []string) error {
	gg := g.generators.Find(args[0])
	if gg == nil {
		return nil
	}

	return gg.Generate(ctx, pwd, args)
}
