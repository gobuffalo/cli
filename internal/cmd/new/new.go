package new

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/cli/internal/genny/assets/standard"
	"github.com/gobuffalo/cli/internal/genny/assets/webpack"
	"github.com/gobuffalo/cli/internal/genny/newapp/api"
	"github.com/gobuffalo/cli/internal/genny/newapp/web"
	"github.com/gobuffalo/cli/internal/genny/vcs"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/gobuffalo/logger"

	pop "github.com/gobuffalo/buffalo-pop/v3/genny/newapp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Creates a new Buffalo application",
	RunE:  RunE,
}

func RunE(cmd *cobra.Command, args []string) error {
	// Restore default values after usage (useful for testing)
	defer func() {
		cmd.Flags().Visit(func(f *pflag.Flag) {
			f.Value.Set(f.DefValue)
		})
		viper.BindPFlags(cmd.Flags())
	}()

	nopts, err := parseNewOptions(args)
	if err != nil {
		return err
	}

	opts := nopts.Options
	app := opts.App

	ctx := context.Background()

	run := genny.WetRunner(ctx)
	lg := logger.New(logger.DebugLevel)
	run.Logger = lg

	if nopts.DryRun {
		run = genny.DryRunner(ctx)
	}

	run.Root = app.Root
	if nopts.Force {
		os.RemoveAll(app.Root)
	}

	var gg *genny.Group

	if app.AsAPI {
		gg, err = api.New(&api.Options{
			Options: opts,
		})
	} else {
		wo := &web.Options{
			Options: opts,
		}
		if app.WithWebpack {
			wo.Webpack = &webpack.Options{}
		} else {
			wo.Standard = &standard.Options{}
		}
		gg, err = web.New(wo)
	}

	if err != nil {
		return err
	}

	g := genny.New()
	g.Command(exec.Command("go", "mod", "tidy"))
	gg.Add(g)

	g = genny.New()
	g.Command(exec.Command("go", "mod", "download"))
	gg.Add(g)

	run.WithGroup(gg)

	if err := run.WithNew(gogen.Fmt(app.Root)); err != nil {
		return err
	}

	// setup VCS last
	if opts.VCS != nil {
		// add the VCS generator
		if err := run.WithNew(vcs.New(opts.VCS)); err != nil {
			return err
		}
	}

	if err := run.Run(); err != nil {
		return err
	}

	run.Logger.Infof("Congratulations! Your application, %s, has been successfully built!", app.Name)
	run.Logger.Infof("You can find your new application at: %v", app.Root)
	run.Logger.Info("Please read the README.md file in your new application for next steps on running your application.")

	return nil
}

func init() {
	Cmd.Flags().Bool("api", false, "skip all front-end code and configure for an API server")
	Cmd.Flags().BoolP("force", "f", false, "delete and remake if the app already exists")
	Cmd.Flags().BoolP("dry-run", "d", false, "dry run")
	Cmd.Flags().BoolP("verbose", "v", false, "verbosely print out the go get commands")

	Cmd.Flags().Bool("skip-pop", false, "skips adding pop/soda to your app")
	Cmd.Flags().Bool("skip-webpack", false, "skips adding Webpack to your app")
	Cmd.Flags().Bool("skip-yarn", false, "use npm instead of yarn for frontend dependencies management")
	Cmd.Flags().Bool("skip-docker", false, "skips generating the Dockerfile")

	Cmd.Flags().String("db-type", "postgres", fmt.Sprintf("specify the type of database you want to use [%s]", strings.Join(pop.AvailableDialects, ", ")))
	Cmd.Flags().String("ci-provider", "none", "specify the type of ci file you would like buffalo to generate [none, travis, gitlab-ci, circleci]")
	Cmd.Flags().String("vcs", "git", "specify the Version control system you would like to use [none, git, bzr]")
	Cmd.Flags().String("module", "", "specify the root module (package) name. [defaults to 'automatic']")

	viper.BindPFlags(Cmd.Flags())

	cfgFile := Cmd.PersistentFlags().String("config", "", "config file (default is $HOME/.buffalo.yaml)")
	skipConfig := Cmd.Flags().Bool("skip-config", false, "skips using the config file")

	cobra.OnInitialize(initConfig(skipConfig, cfgFile))
}
