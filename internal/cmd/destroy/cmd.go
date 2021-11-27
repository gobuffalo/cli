package destroy

import (
	"github.com/spf13/cobra"
)

// DestroyCmd destroys generated code
var cmd = &cobra.Command{
	Use:     "destroy",
	Short:   "Destroy generated components",
	Aliases: []string{"d"},
}

func Cmd() *cobra.Command {
	cmd.AddCommand(resourceCmd)
	cmd.AddCommand(actionCmd)
	cmd.AddCommand(mailerCmd)

	cmd.PersistentFlags().BoolVarP(&yesToAll, "yes", "y", false, "confirms all beforehand")

	return cmd
}
