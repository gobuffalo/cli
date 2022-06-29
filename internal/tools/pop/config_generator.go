package pop

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/cli/internal/defaults"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/pop/v6/genny/config"
)

var ConfigGenerator = &configGenerator{}

type configGenerator struct {
	configFile string
	dialect    string // default: postgres
}

func (c configGenerator) Name() string {
	return "pop/config"
}

func (c configGenerator) HelpText() string {
	return "Generates a database.yml file for your project."
}

func (c *configGenerator) Generate(ctx context.Context, pwd string, args []string) error {
	//defaults.String(cflagVal, "database.yml")

	cfgFile := defaults.String(c.configFile, "database.yml")
	run := genny.WetRunner(ctx)

	g, err := config.New(&config.Options{
		Root:     pwd,
		Prefix:   filepath.Base(pwd),
		FileName: cfgFile,
		Dialect:  c.dialect,
	})

	if err != nil {
		return err
	}

	run.With(g)

	return run.Run()
}

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"strings"

// 	"github.com/gobuffalo/genny/v2"
// 	"github.com/gobuffalo/pop/v6"
// 	"github.com/gobuffalo/pop/v6/genny/config"
// 	"github.com/gobuffalo/pop/v6/internal/defaults"
// 	"github.com/spf13/cobra"
// )

// func init() {
// 	ConfigCmd.Flags().StringVarP(&dialect, "type", "t", "postgres", fmt.Sprintf("The type of database you want to use (%s)", strings.Join(pop.AvailableDialects, ", ")))
// }

// var dialect string

// // ConfigCmd is the command to generate pop config files
// var ConfigCmd = &cobra.Command{
// 	Use:              "config",
// 	Short:            "Generates a database.yml file for your project.",
// 	PersistentPreRun: func(c *cobra.Command, args []string) {},
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		cflag := cmd.Flag("config")
// 		cflagVal := ""
// 		if cflag != nil {
// 			cflagVal = cflag.Value.String()
// 		}

// 	},
// }
