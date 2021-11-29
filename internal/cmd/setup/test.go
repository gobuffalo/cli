package setup

import (
	"fmt"
	"os/exec"

	"github.com/gobuffalo/meta"
)

func testCheck(meta.App) error {
	err := run(exec.Command("buffalo", "test"))
	if err != nil {
		return fmt.Errorf("We encountered the following error when trying to run your applications tests:\n%s", err)
	}
	return nil
}
