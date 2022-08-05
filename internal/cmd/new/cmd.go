package new

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/pop/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [name]",
		Short: "Creates a new Buffalo application",
		RunE:  RunE,
	}

	cmd.Flags().Bool("api", false, "skip all front-end code and configure for an API server")
	cmd.Flags().BoolP("force", "f", false, "delete and remake if the app already exists")
	cmd.Flags().BoolP("dry-run", "d", false, "dry run")
	cmd.Flags().BoolP("verbose", "v", false, "verbosely print out the go get commands")
	cmd.Flags().Bool("skip-pop", false, "skips adding pop/soda to your app")
	cmd.Flags().Bool("skip-webpack", false, "skips adding Webpack to your app")
	cmd.Flags().Bool("skip-yarn", false, "use npm instead of yarn for frontend dependencies management")
	cmd.Flags().Bool("skip-docker", false, "skips generating the Dockerfile")
	cmd.Flags().String("db-type", "postgres", fmt.Sprintf("specify the type of database you want to use [%s]", strings.Join(pop.AvailableDialects, ", ")))
	cmd.Flags().String("ci-provider", "none", "specify the type of ci file you would like buffalo to generate [none, circleci, github, gitlab-ci, travis]")
	cmd.Flags().String("vcs", "git", "specify the Version control system you would like to use [none, git, bzr]")
	cmd.Flags().String("module", "", "specify the root module (package) name. [defaults to 'automatic']")

	viper.BindPFlags(cmd.Flags())

	cfgFile := cmd.PersistentFlags().String("config", "", "config file (default is $HOME/.buffalo.yaml)")
	skipConfig := cmd.Flags().Bool("skip-config", false, "skips using the config file")
	cobra.OnInitialize(initConfig(skipConfig, cfgFile))

	return cmd
}
