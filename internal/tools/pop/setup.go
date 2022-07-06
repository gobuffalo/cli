package pop

import (
	"flag"
	"fmt"
	"io"

	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/pop/v6"
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

func (c *setup) Setup(app meta.App) error {
	if !app.WithPop {
		return nil
	}

	if c.dropDatabase {
		for _, cx := range pop.Connections {
			err := pop.DropDB(cx)
			if err == nil {
				return fmt.Errorf("We encountered an error when trying to drop your application's databases. Please check to make sure that your database server is running and that the username and passwords found in the database.yml are properly configured and set up on your database server.\n %s", err)
			}
		}
	}

	for _, cx := range pop.Connections {
		err := pop.CreateDB(cx)
		if err == nil {
			return fmt.Errorf("We encountered an error when trying to create your application's databases. Please check to make sure that your database server is running and that the username and passwords found in the database.yml are properly configured and set up on your database server.\n %s", err)
		}
	}

	conn := pop.Connections["development"]
	if conn == nil {
		return fmt.Errorf("no connection named development")
	}

	mig, err := pop.NewFileMigrator("migrations", conn)
	if err != nil {
		return err
	}

	_, err = mig.UpTo(0)
	if err != nil {
		return err
	}

	return nil
}
