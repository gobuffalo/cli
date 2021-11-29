package version

import (
	"encoding/json"
	"os"

	"github.com/gobuffalo/cli/internal/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func run(c *cobra.Command, args []string) {
	if !jsonOutput {
		logrus.Infof("Buffalo version is: %s", runtime.Version)
		return
	}

	build := runtime.BuildInfo{}
	build.Version = runtime.Version

	enc := json.NewEncoder(os.Stderr)
	enc.SetIndent("", "    ")
	enc.Encode(build)
}
