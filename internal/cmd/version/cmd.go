package version

import "github.com/spf13/cobra"

var jsonOutput bool

func Cmd() *cobra.Command {
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Print information in json format")

	return cmd
}
