package cmd

import "github.com/gobuffalo/cli/internal/cmd/new"

func init() {
	decorate("new", new.Cmd)

	RootCmd.AddCommand(new.Cmd)
}
