package setup

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/gobuffalo/cli/cmd/cli/help"
	"github.com/gobuffalo/cli/cmd/cli/plugin"
	"github.com/gobuffalo/events"
	"github.com/gobuffalo/meta"
)

const (
	// EvtSetupStarted is emitted when `buffalo setup` starts
	EvtSetupStarted = "buffalo:setup:started"
	// EvtSetupErr is emitted if `buffalo setup` fails
	EvtSetupErr = "buffalo:setup:err"
	// EvtSetupFinished is emitted when `buffalo setup` finishes
	EvtSetupFinished = "buffalo:setup:finished"
)

var (
	//go:embed setupdescription.txt
	setupLongDescription string

	// Command instance of the setup command to be used outside of this package.
	Command = &command{}
)

type command struct {
	setupers []Setuper
}

func (c command) Name() string {
	return "setup"
}

func (c command) HelpText() string {
	return "Setups a newly created, or recently checked out application."
}

func (c *command) LongHelpText() string {
	s := "Registered Setupers\n"
	for _, v := range c.setupers {
		description := ""
		if st, ok := v.(help.HelpTexter); ok {
			description = st.HelpText()
		}

		s += fmt.Sprintf("%s\t\t%v\n", v.Name(), description)
	}

	return s
}

func (c *command) Main(ctx context.Context, pwd string, args []string) error {
	app := meta.New(".")
	payload := events.Payload{
		"app": app,
	}
	events.EmitPayload(EvtSetupStarted, payload)

	for _, st := range c.setupers {
		err := st.Setup(app)
		if err == nil {
			continue
		}

		// Emit an error if something went wrong
		events.EmitError(EvtSetupErr, err, payload)
		return err
	}

	events.EmitPayload(EvtSetupFinished, payload)
	return nil
}

func (c *command) Receive(plugins plugin.Plugins) {
	for _, v := range plugins {
		v, ok := v.(Setuper)
		if !ok {
			continue
		}

		c.setupers = append(c.setupers, v)
	}
}
