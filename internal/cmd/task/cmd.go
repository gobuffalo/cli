package task

import (
	grifts "github.com/markbates/grift/cmd"
	"github.com/spf13/cobra"
)

// task command is a forward to grift tasks
var cmd = &cobra.Command{
	Use:                "task",
	Aliases:            []string{"t", "tasks"},
	Short:              "Run grift tasks",
	DisableFlagParsing: true,
	RunE: func(c *cobra.Command, args []string) error {
		return grifts.Run("buffalo task", args)
	},
}

func Cmd() *cobra.Command {
	return cmd
}
