package task

import (
	grifts "github.com/markbates/grift/cmd"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "task",
		Aliases:            []string{"t", "tasks"},
		Short:              "Run grift tasks",
		DisableFlagParsing: true,
		RunE: func(c *cobra.Command, args []string) error {
			return grifts.Run("buffalo task", args)
		},
	}

	return cmd
}
