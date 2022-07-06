package pop

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/gobuffalo/pop/v6"
)

var Migrate = &migrate{}

type migrate struct {
	flagSet *flag.FlagSet

	steps int
	env   string
}

func (c migrate) Name() string {
	return "migrate"
}

func (c migrate) HelpText() string {
	return "Runs migrations up or down. Also, provides the status of the migrations."
}

func (c migrate) Usage() string {
	return "buffalo db migrate [flags] [up|down|status] "
}

func (c *migrate) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(io.Discard)
	}

	defaultSteps := 0
	if len(args) > 0 && strings.Contains(strings.Join(args, ","), "down") {
		defaultSteps = 1
	}

	c.flagSet.IntVar(&c.steps, "steps", defaultSteps, "number of steps to migrate")
	c.flagSet.StringVar(&c.env, "env", "development", "environment or connection name to migrate")

	fmt.Println("migrate:", args)

	_ = c.flagSet.Parse(args[1:])

	return c.flagSet, nil
}

func (c migrate) PopMain(ctx context.Context, pwd string, args []string) error {
	action := "up"
	if len(args) > 0 {
		action = args[0]
	}

	fmt.Println("steps", c.steps)

	conn := pop.Connections[c.env]
	if conn == nil {
		return fmt.Errorf("no connection named %s", c.env)
	}

	mig, err := pop.NewFileMigrator("migrations", conn)
	if err != nil {
		return err
	}

	switch action {
	case "up":
		_, err = mig.UpTo(c.steps)
		return err
	case "down":
		return mig.Down(c.steps)
	case "status":
		return mig.Status(os.Stdout)
	}

	return nil
}
