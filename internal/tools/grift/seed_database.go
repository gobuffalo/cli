package grift

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/meta"
)

var SetupSeedDatabase = setupSeedDatabase{}

type setupSeedDatabase struct{}

func (c setupSeedDatabase) Name() string {
	return "grift/seed"
}

func (c setupSeedDatabase) HelpText() string {
	return "Attempts to seed the database with the `seed` Grift task."
}

func (c setupSeedDatabase) Setup(app meta.App) error {
	// Trying to seed the database with the `seed task`
	cmd := exec.Command("buffalo", "t", "list")
	out, err := cmd.Output()
	if err != nil {
		// no tasks configured, so return
		return nil
	}

	if bytes.Contains(out, []byte("db:seed")) {
		err := run(exec.Command("buffalo", "task", "db:seed"))
		if err != nil {
			return fmt.Errorf("We encountered the following error when trying to seed your database:\n%s", err)
		}
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
