package setup

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup a newly created, or recently checked out application.",
		Long:  setupLongDescription,
		RunE:  runE,
	}

	cmd.Flags().BoolVarP(&setupOptions.verbose, "verbose", "v", false, "run with verbose output")
	cmd.Flags().BoolVarP(&setupOptions.dropDatabases, "drop", "d", false, "drop existing databases")

	return cmd
}
