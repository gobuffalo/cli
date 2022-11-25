package test

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "test",
	Short: "Runs tests for your Buffalo app",
	Long:  ``,

	// DisableFlagParsing is set to true since we will need to allow undefined
	// flags to be passed to the go test command.
	DisableFlagParsing: true,

	RunE: func(c *cobra.Command, args []string) error {
		// Set the environment to be test so that the rest of the tooling
		// understands we're in testing mode.
		if err := os.Setenv("GO_ENV", "test"); err != nil {
			return err
		}

		// Setup the database before running the tests.
		if err := setupDatabase(args); err != nil {
			return err
		}

		tcmd, err := buildCmd(args)
		if err != nil {
			return err
		}

		fmt.Println("[INFO]", strings.Join(tcmd.Args, " "))
		return tcmd.Run()
	},
}

// TODO: Remove me
func Cmd() *cobra.Command {
	return cmd
}
