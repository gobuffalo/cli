package new

import (
	"context"
	"os"
	"os/exec"

	"github.com/gobuffalo/cli/internal/genny/assets/standard"
	"github.com/gobuffalo/cli/internal/genny/assets/webpack"
	"github.com/gobuffalo/cli/internal/genny/newapp/api"
	"github.com/gobuffalo/cli/internal/genny/newapp/web"
	"github.com/gobuffalo/cli/internal/genny/vcs"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/gobuffalo/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

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
		if gg, err = api.New(&api.Options{
			Options: opts,
		}); err != nil {
			return err
		}
	} else {
		wo := &web.Options{
			Options: opts,
		}
		if app.WithWebpack {
			wo.Webpack = &webpack.Options{}
		} else {
			wo.Standard = &standard.Options{}
		}
		if gg, err = web.New(wo); err != nil {
			return err
		}
	}

	g := genny.New()
	g.Command(exec.Command("go", "mod", "tidy"))
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

	run.Logger.Infof("Congratulations! Your application, %s, has been successfully generated!", app.Name)
	run.Logger.Infof("You can find your new application at: %v", app.Root)
	run.Logger.Info("Please read the README.md file in your new application for next steps on running your application.")

	return nil
}
