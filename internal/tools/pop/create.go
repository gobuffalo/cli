package pop

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/gobuffalo/pop/v6"
)

var Create = &create{}

type create struct {
	flagSet *flag.FlagSet

	all bool
	env string
}

func (c create) Name() string {
	return "create"
}

func (c create) HelpText() string {
	return "Creates databases for you"
}

func (c create) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name(), flag.ExitOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(io.Discard)

		c.flagSet.BoolVar(&c.all, "all", false, "create all databases")
		c.flagSet.StringVar(&c.env, "env", "development", "environment or connection name to create")
	}

	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}

func (c create) PopMain(ctx context.Context, pwd string, args []string) error {
	var err error
	if !c.all {
		conn := pop.Connections[c.env]
		if conn == nil {
			return fmt.Errorf("no connection named %s", c.env)
		}

		return pop.CreateDB(conn)
	}

	for _, cx := range pop.Connections {
		err := pop.CreateDB(cx)
		if err == nil {
			return err
		}
	}

	return err
}
