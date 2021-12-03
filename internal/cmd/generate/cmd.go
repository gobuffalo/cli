package generate

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generate application components",
		Aliases: []string{"g"},
	}

	cmd.AddCommand(ResourceCmd)
	cmd.AddCommand(ActionCmd)
	cmd.AddCommand(TaskCmd)
	cmd.AddCommand(MailCmd)

	return cmd
}
