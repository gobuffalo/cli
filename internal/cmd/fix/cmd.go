package fix

import (
	"fmt"

	"github.com/gobuffalo/cli/internal/runtime"
	"github.com/spf13/cobra"
)

// yesToAll will be used by the command to skip the confirmation
// and perform all implied destroy operations
var yesToAll bool = false

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fix",
		Aliases: []string{"update"},
		Short:   fmt.Sprintf("Attempt to fix a Buffalo application's API to match version %s", runtime.Version),
		RunE:    RunE,
	}

	cmd.Flags().BoolVarP(&yesToAll, "y", "y", false, "update all without asking for confirmation")

	return cmd
}
