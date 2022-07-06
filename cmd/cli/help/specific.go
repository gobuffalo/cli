package help

import (
	"flag"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/gobuffalo/cli/cmd/cli/clio"
	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

// Prints the help for the given plugin
func Specific(stdout io.Writer, cm plugin.Plugin) error {
	usage := fmt.Sprintf("buffalo %v [options]", cm.Name())
	if fp, ok := cm.(Usager); ok {
		usage = fp.Usage()
	}

	fmt.Fprintf(stdout, "Usage: %v \n\n", usage)

	if ht, ok := cm.(HelpTexter); ok {
		fmt.Fprintf(stdout, ht.HelpText()+"\n\n")
	}

	if ht, ok := cm.(plugin.Aliaser); ok {
		fmt.Fprintf(stdout, "Aliases:\n")
		fmt.Fprintf(stdout, "  %s\n", strings.Join(ht.Aliases(), ", "))
		fmt.Fprintf(stdout, "\n")
	}

	if ht, ok := cm.(LongHelpTexter); ok {
		fmt.Fprintf(stdout, ht.LongHelpText()+"\n\n")
	}

	if fl, ok := cm.(clio.FlagParser); ok {
		fx, _ := fl.ParseFlags([]string{"xxx"})
		w := tabwriter.NewWriter(stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintf(w, "Flags:\n")

		fx.VisitAll(func(ff *flag.Flag) {
			fmt.Fprintf(w, "  --%v\t%v\n", ff.Name, ff.Usage)
		})

		w.Flush()
	}

	return nil
}
