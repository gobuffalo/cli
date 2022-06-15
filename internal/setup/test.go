package setup

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/meta"
)

const Test = test("test-setup")

// Test type is a Setupper that will run the tests
// when the setup command is invoked.
type test string

func (ts test) Name() string {
	return "setup/test"
}

func (ts test) HelpText() string {
	return "Runs the application tests"
}

func (ts test) Setup(app meta.App) error {

	var run = func(cmd *exec.Cmd) error {
		fmt.Println("--> %s", strings.Join(cmd.Args, " "))
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return cmd.Run()
	}

	err := run(exec.Command("buffalo", "test"))
	if err != nil {
		return fmt.Errorf("We encountered the following error when trying to run your applications tests:\n%s", err)
	}

	return nil
}
