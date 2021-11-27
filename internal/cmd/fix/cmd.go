package fix

import (
	"fmt"

	"github.com/gobuffalo/cli/internal/runtime"
	"github.com/spf13/cobra"
)

// cmd represents the info command
var cmd = &cobra.Command{
	Use:     "fix",
	Aliases: []string{"update"},
	Short:   fmt.Sprintf("Attempt to fix a Buffalo application's API to match version %s", runtime.Version),
	RunE: func(cmd *cobra.Command, args []string) error {
		return Run()
	},
}

func Cmd() *cobra.Command {
	cmd.Flags().BoolVarP(&YesToAll, "y", "y", false, "update all without asking for confirmation")

	return cmd
}
