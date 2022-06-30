package pop

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/meta"
)

var Setup = &setup{
	flagSet: flag.NewFlagSet("setup", flag.ExitOnError),
}

type setup struct {
	flagSet *flag.FlagSet

	verbose      bool
	dropDatabase bool
}

func (c setup) Name() string {
	return "pop/setup"
}

func (c setup) HelpText() string {
	return "Setups the database"
}

func (c *setup) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet("setup", flag.ExitOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(io.Discard)
	}

	c.flagSet.BoolVar(&c.verbose, "verbose", false, "run with verbose output")
	c.flagSet.BoolVar(&c.dropDatabase, "drop", false, "drop existing databases")

	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}

// TODO: Move these exec to use the package
func (c *setup) Setup(app meta.App) error {
	if !app.WithPop {
		return nil
	}

	if c.dropDatabase {
		err := run(exec.Command("buffalo", "pop", "drop", "-a"))
		if err != nil {
			return fmt.Errorf("We encountered an error when trying to drop your application's databases. Please check to make sure that your database server is running and that the username and passwords found in the database.yml are properly configured and set up on your database server.\n %s", err)
		}
	}

	err := run(exec.Command("buffalo", "pop", "create", "-a"))
	if err != nil {
		return fmt.Errorf("We encountered an error when trying to create your application's databases. Please check to make sure that your database server is running and that the username and passwords found in the database.yml are properly configured and set up on your database server.\n %s", err)
	}

	// Running migrations on the database with Pop
	err = run(exec.Command("buffalo", "pop", "migrate"))
	if err != nil {
		return fmt.Errorf("We encountered the following error when trying to migrate your database:\n%s", err)
	}

	return nil

}

func run(cmd *exec.Cmd) error {
	fmt.Printf("--> %s\n", strings.Join(cmd.Args, " "))

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
