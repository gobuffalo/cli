package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/gobuffalo/cli/internal/cmd/build"
	"github.com/gobuffalo/cli/internal/cmd/destroy"
	"github.com/gobuffalo/cli/internal/cmd/dev"
	"github.com/gobuffalo/cli/internal/cmd/fix"
	"github.com/gobuffalo/cli/internal/cmd/generate"
	"github.com/gobuffalo/cli/internal/cmd/info"
	"github.com/gobuffalo/cli/internal/cmd/new"
	cmdplugins "github.com/gobuffalo/cli/internal/cmd/plugins"
	"github.com/gobuffalo/cli/internal/cmd/routes"
	"github.com/gobuffalo/cli/internal/cmd/setup"
	"github.com/gobuffalo/cli/internal/cmd/task"
	"github.com/gobuffalo/cli/internal/cmd/test"
	"github.com/gobuffalo/cli/internal/cmd/version"
	"github.com/gobuffalo/cli/internal/plugins"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	dbNotFound  = `unknown command "db"`
	popNotFound = `unknown command "pop"`
)

var (
	anywhereCommands = []string{"new", "version", "info", "help"}

	//go:embed popinstructions.txt
	popInstallInstructions string
)

// rootCmd is the hook for all of the other commands in the buffalo binary.
var rootCmd = &cobra.Command{
	SilenceErrors: true,
	Use:           "buffalo",
	Short:         "Build Buffalo applications with ease",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := plugins.Load(); err != nil {
			return err
		}

		isFreeCommand := false
		for _, freeCmd := range anywhereCommands {
			if freeCmd != cmd.Name() {
				continue
			}

			isFreeCommand = true
			break
		}

		if isFreeCommand {
			return nil
		}

		if !insideBuffaloProject() {
			return fmt.Errorf("you need to be inside your buffalo project path to run this command")
		}

		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	newCmd := new.Cmd()
	setupCmd := setup.Cmd()
	generateCmd := generate.Cmd()
	destroyCmd := destroy.Cmd()
	versionCmd := version.Cmd()
	testCmd := test.Cmd()
	devCmd := dev.Cmd()
	taskCmd := task.Cmd()
	routesCmd := routes.Cmd()

	decorate("new", newCmd)
	decorate("info", rootCmd)
	decorate("fix", rootCmd)
	decorate("update", rootCmd)
	decorate("setup", setupCmd)
	decorate("generate", generateCmd)
	decorate("destroy", destroyCmd)
	decorate("version", versionCmd)
	decorate("test", testCmd)
	decorate("dev", devCmd)
	decorate("task", taskCmd)
	decorate("routes", routesCmd)

	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(build.Cmd())
	rootCmd.AddCommand(cmdplugins.PluginsCmd)
	rootCmd.AddCommand(info.Cmd())
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(fix.Cmd())
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(taskCmd)
	rootCmd.AddCommand(routesCmd)

	decorate("root", rootCmd)
	err := rootCmd.Execute()
	if err == nil {
		return
	}

	if strings.Contains(err.Error(), dbNotFound) || strings.Contains(err.Error(), popNotFound) {
		logrus.Errorf(popInstallInstructions)
		os.Exit(-1)
	}
	logrus.Errorf("Error: %s", err)
	if strings.Contains(err.Error(), dbNotFound) || strings.Contains(err.Error(), popNotFound) {
		fmt.Println(popInstallInstructions)
		os.Exit(-1)
	}
	os.Exit(-1)
}

func insideBuffaloProject() bool {
	if _, err := os.Stat(".buffalo.dev.yml"); err != nil {
		return false
	}

	return true
}
