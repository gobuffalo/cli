package build

import (
	"io"

	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

// Version runner will run a version cmd with a given
// io.Writer. This is useful to extract the verion number
// when building.
type VersionRunner interface {
	plugin.Plugin

	RunVersionCmd(io.Writer) error
}
