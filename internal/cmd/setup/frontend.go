package setup

import (
	"fmt"
	"os/exec"

	"github.com/gobuffalo/meta"
)

func yarnCheck(app meta.App) error {
	if err := nodeCheck(app); err != nil {
		return err
	}
	if _, err := exec.LookPath("yarnpkg"); err != nil {
		err := run(exec.Command("npm", "install", "-g", "yarn@berry"))
		if err != nil {
			return fmt.Errorf("This application require yarn, and we could not find it installed on your system. We tried to install it for you, but ran into the following error:\n%s", err)
		}
	}
	if err := run(exec.Command("yarnpkg", "install", "--silent")); err != nil {
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

func assetCheck(app meta.App) error {
	if !app.WithWebpack {
		return nil
	}
	if app.WithYarn {
		return yarnCheck(app)
	}
	return npmCheck(app)
}

func npmCheck(app meta.App) error {
	err := nodeCheck(app)
	if err != nil {
		return err
	}
	err = run(exec.Command("npm", "install", "--no-progress"))
	if err != nil {
		return fmt.Errorf("We encountered the following error when trying to install your asset dependencies using npm:\n%s", err)
	}
	return nil
}
