package pop

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/gobuffalo/pop/v6"
)

var Drop = &drop{}

type drop struct {
	flagSet *flag.FlagSet

	all bool
	env string
}

func (c drop) Name() string {
	return "drop"
}

func (c drop) Usage() string {
	return "buffalo db drop [flags]"
}

func (c drop) HelpText() string {
	return "Drops databases for you"
}

func (c drop) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name(), flag.ExitOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(io.Discard)

		c.flagSet.BoolVar(&c.all, "all", false, "create all databases")
		c.flagSet.StringVar(&c.env, "env", "development", "environment or connection name to drop")
	}

	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}

func (c drop) PopMain(ctx context.Context, pwd string, args []string) error {
	if !c.all {
		conn := pop.Connections[c.env]
		if conn == nil {
			return fmt.Errorf("no connection named %s", c.env)
		}

		return pop.CreateDB(conn)
	}

	for _, conn := range pop.Connections {
		err := pop.DropDB(conn)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
