package destroy

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "destroy",
		Short:   "Destroy generated components",
		Aliases: []string{"d"},
	}

	cmd.AddCommand(resourceCmd)
	cmd.AddCommand(actionCmd)
	cmd.AddCommand(mailerCmd)

	cmd.PersistentFlags().BoolVarP(&yesToAll, "yes", "y", false, "confirms all beforehand")

	return cmd
}
