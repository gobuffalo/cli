package version

import (
	"encoding/json"
	"os"

	"github.com/gobuffalo/cli/internal/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Run: func(c *cobra.Command, args []string) {
		if !jsonOutput {
			logrus.Infof("Buffalo version is: %s", runtime.Version)
			return
		}

		build := runtime.BuildInfo{}
		build.Version = runtime.Version

		enc := json.NewEncoder(os.Stderr)
		enc.SetIndent("", "    ")
		enc.Encode(build)
	},
	// needed to override the root level pre-run func
	PersistentPreRunE: func(c *cobra.Command, args []string) error {
		return nil
	},
}
