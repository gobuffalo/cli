package frontend

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/meta"
)

var Setup = &setup{}

type setup struct{}

func (sf setup) Name() string {
	return "frontend/setup"
}

func (sf setup) HelpText() string {
	return "Runs yarn or npm install depending on the application's configuration."
}

func (sf setup) Setup(app meta.App) error {
	if !app.WithWebpack {
		return nil
	}

	err := nodeCheck(app)
	if err != nil {
		return err
	}

	if app.WithYarn {
		return yarnCheck(app)
	}

	return npmCheck(app)
}

func yarnCheck(app meta.App) error {
	if _, err := exec.LookPath("yarnpkg"); err != nil {
		err := run(exec.Command("npm", "install", "-g", "yarn"))
		if err != nil {
			return fmt.Errorf("This application require yarn, and we could not find it installed on your system. We tried to install it for you, but ran into the following error:\n%s", err)
		}
	}

	if err := run(exec.Command("yarnpkg", "install")); err != nil {
		return fmt.Errorf("We encountered the following error when trying to install your asset dependencies using yarn:\n%s", err)
	}

	return nil
}

func nodeCheck(meta.App) error {
	if _, err := exec.LookPath("node"); err != nil {
		return fmt.Errorf("this application requires node, and we could not find it installed on your system please install node and try again")
	}

	if _, err := exec.LookPath("npm"); err != nil {
		return fmt.Errorf("this application requires npm, and we could not find it installed on your system please install npm and try again")
	}

	return nil
}

func npmCheck(app meta.App) error {
	err := run(exec.Command("npm", "install", "--no-progress"))
	if err == nil {
		return nil
	}

	return fmt.Errorf("We encountered the following error when trying to install your asset dependencies using npm:\n%s", err)
}

func run(cmd *exec.Cmd) error {
	fmt.Printf("--> %s\n", strings.Join(cmd.Args, " "))

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
