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
	cmd.AddCommand(ResourceCmd)
	cmd.AddCommand(ActionCmd)
	cmd.AddCommand(MailerCmd)

	cmd.PersistentFlags().BoolVarP(&YesToAll, "yes", "y", false, "confirms all beforehand")

	return cmd
}
