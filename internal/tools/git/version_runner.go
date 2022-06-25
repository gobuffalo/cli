package git

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

var VersionRunner = gitVersionRunner("git")

type gitVersionRunner string

func (gv gitVersionRunner) Name() string {
	return string(gv)
}

func (gv gitVersionRunner) RunVersionCmd(out io.Writer) error {
	// If .git folder does not exist return default version
	if stat, err := os.Stat(".git"); err != nil || !stat.IsDir() {
		return fmt.Errorf("could not find .git folder")
	}

	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Stdout = out

	return cmd.Run()
}
