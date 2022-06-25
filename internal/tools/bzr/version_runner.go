package bzr

import (
	"io"
	"os/exec"
)

var VersionRunner = bzrVersionRunner("bzr")

type bzrVersionRunner string

func (gv bzrVersionRunner) Name() string {
	return string(gv)
}

func (gv bzrVersionRunner) RunVersionCmd(out io.Writer) error {
	cmd := exec.Command("bzr", "revno")
	cmd.Stdout = out

	return cmd.Run()
}
