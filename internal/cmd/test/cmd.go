package test

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "test",
		Short:              "Run the tests for the Buffalo app. Use --force-migrations to skip schema load.",
		DisableFlagParsing: true,
		RunE:               runE,
	}

	return cmd
}
