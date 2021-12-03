package version

import "github.com/spf13/cobra"

// jsonOutput for the version command
var jsonOutput bool = false

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run:   run,
		// needed to override the root level pre-run func
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Print information in json format")

	return cmd
}
