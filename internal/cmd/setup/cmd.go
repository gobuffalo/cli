package setup

import (
	_ "embed"

	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd.Flags().BoolVarP(&setupOptions.verbose, "verbose", "v", false, "run with verbose output")
	cmd.Flags().BoolVarP(&setupOptions.dropDatabases, "drop", "d", false, "drop existing databases")

	return cmd
}
