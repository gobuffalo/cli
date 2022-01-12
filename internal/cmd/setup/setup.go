package setup

import (
	_ "embed"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/events"
	"github.com/gobuffalo/meta"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type setupCheck func(meta.App) error

const (
	// EvtSetupStarted is emitted when `buffalo setup` starts
	EvtSetupStarted = "buffalo:setup:started"
	// EvtSetupErr is emitted if `buffalo setup` fails
	EvtSetupErr = "buffalo:setup:err"
	// EvtSetupFinished is emitted when `buffalo setup` finishes
	EvtSetupFinished = "buffalo:setup:finished"
)

var (
	setupOptions = struct {
		verbose       bool
		dropDatabases bool
	}{}

	//go:embed setupdescription.txt
	setupLongDescription string

	// checks that will be run by the setup command
	checks = []setupCheck{
		assetCheck,
		databaseCheck,
		testCheck,
	}
)

func runE(cmd *cobra.Command, args []string) error {
	app := meta.New(".")
	payload := events.Payload{
		"app": app,
	}
	events.EmitPayload(EvtSetupStarted, payload)

	for _, check := range checks {
		err := check(app)
		if err != nil {
			events.EmitError(EvtSetupErr, err, payload)
			return err
		}
	}

	events.EmitPayload(EvtSetupFinished, payload)
	return nil
}

func run(cmd *exec.Cmd) error {
	logrus.Infof("--> %s", strings.Join(cmd.Args, " "))
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
