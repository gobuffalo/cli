package pop

import (
	"context"
	"flag"
	"io"

	"github.com/gobuffalo/pop/v6"
)

var Create = &create{}

type create struct {
	flagSet *flag.FlagSet
	all     bool
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
	}

	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}

func (c create) Run(ctx context.Context, conn *pop.Connection) error {
	var err error
	if !c.all {
		return pop.CreateDB(conn)
	}

	for _, conn := range pop.Connections {
		err := pop.CreateDB(conn)
		if err == nil {
			return err
		}
	}

	return err
}
