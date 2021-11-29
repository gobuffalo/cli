package routes

import (
	grifts "github.com/markbates/grift/cmd"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routes",
		Short: "Print all defined routes",
		RunE: func(c *cobra.Command, args []string) error {
			return grifts.Run("buffalo task", []string{"routes"})
		},
	}

	return cmd
}
