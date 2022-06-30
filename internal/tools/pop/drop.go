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
	all     bool
}

func (c drop) Name() string {
	return "drop"
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
	}

	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}

func (c drop) Run(ctx context.Context, conn *pop.Connection) error {
	if !c.all {
		return pop.DropDB(conn)
	}

	for _, conn := range pop.Connections {
		err := pop.DropDB(conn)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
