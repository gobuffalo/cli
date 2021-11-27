package generate

import (
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate application components",
	Aliases: []string{"g"},
}

func Cmd() *cobra.Command {
	cmd.AddCommand(ResourceCmd)
	cmd.AddCommand(ActionCmd)
	cmd.AddCommand(TaskCmd)
	cmd.AddCommand(MailCmd)

	return cmd
}
