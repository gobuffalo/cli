package new

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	pop "github.com/gobuffalo/buffalo-pop/v3/genny/newapp"
	"github.com/gobuffalo/cli/internal/genny/assets/standard"
	"github.com/gobuffalo/cli/internal/genny/assets/webpack"
	"github.com/gobuffalo/cli/internal/genny/ci"
	"github.com/gobuffalo/cli/internal/genny/docker"
	"github.com/gobuffalo/cli/internal/genny/newapp/api"
	"github.com/gobuffalo/cli/internal/genny/newapp/core"
	"github.com/gobuffalo/cli/internal/genny/newapp/web"
	"github.com/gobuffalo/cli/internal/genny/refresh"
	"github.com/gobuffalo/cli/internal/genny/vcs"
	fname "github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/meta"
	plib "github.com/gobuffalo/pop/v6"
)

var Command = &command{
	options: core.NewOptions(),
}

type command struct {
	flagSet *flag.FlagSet

	options *core.Options

	Module  string
	Force   bool
	Verbose bool
	DryRun  bool
}

func (c command) Name() string {
	return "new"
}

func (c command) Usage() string {
	return "new [name]"
}

func (c command) HelpText() string {
	return "Creates a new Buffalo application"
}

func (c *command) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet("new", flag.ContinueOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(ioutil.Discard)
	}

	// Adding these here since are negation of
	// the options for application creation.
	var skipPop, skipWebpack, skipYarn, skipDocker bool

	c.flagSet.BoolVar(&c.options.App.AsAPI, "api", false, "skip all front-end code and configure for an API server")
	c.flagSet.BoolVar(&skipPop, "skip-pop", false, "skip all back-end code and configure for a web server")
	c.flagSet.BoolVar(&skipWebpack, "skip-webpack", false, "skip all front-end code and configure for a web server")
	c.flagSet.BoolVar(&skipYarn, "skip-yarn", false, "skip all front-end code and configure for a web server")
	c.flagSet.BoolVar(&skipDocker, "skip-docker", false, "skip all front-end code and configure for a web server")

	c.flagSet.BoolVar(&c.Force, "force", false, "delete and remake if the app already exists")
	c.flagSet.BoolVar(&c.DryRun, "dry-run", false, "dry run")
	c.flagSet.BoolVar(&c.Verbose, "verbose", false, "verbosely print out the go get commands")

	c.flagSet.StringVar(&c.Module, "module", "", "module to use for the application")
	c.flagSet.StringVar(&c.options.App.VCS, "vcs", "git", "specify the Version control system you would like to use [none, git, bzr]")
	c.flagSet.StringVar(&c.options.CI.DBType, "db", "postgres", fmt.Sprintf("specify the type of database you want to use [%s]", strings.Join(plib.AvailableDialects, ", ")))
	c.flagSet.StringVar(&c.options.CI.Provider, "ci-provider", "none", "specify the CI provider you want to use [none, travis, gitlab-ci, circleci]")

	_ = c.flagSet.Parse(args)

	c.options.App.WithDocker = !skipDocker
	c.options.App.WithWebpack = !skipWebpack
	c.options.App.WithYarn = !skipYarn
	c.options.App.WithPop = !skipPop

	c.options.App.WithGrifts = true
	c.options.App.WithNodeJs = c.options.App.WithWebpack
	c.options.App.AsWeb = !c.options.App.AsAPI

	c.options.Refresh = &refresh.Options{}

	if c.options.App.AsAPI {
		c.options.App.WithWebpack = false
		c.options.App.WithYarn = false
		c.options.App.WithNodeJs = false
	}

	if c.options.App.WithDocker {
		c.options.Docker = &docker.Options{}
	}

	if pr := c.options.CI.Provider; pr != "none" {
		c.options.CI = &ci.Options{
			Provider: pr,
			DBType:   c.options.CI.DBType,
		}
	}

	if len(c.options.App.VCS) > 0 && c.options.App.VCS != "none" {
		c.options.VCS = &vcs.Options{
			Provider: c.options.App.VCS,
		}
	}

	return c.flagSet, nil
}

func (c *command) Main(ctx context.Context, pwd string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must enter a name for your new application")
	}

	c.options.App = meta.New(pwd)
	c.options.App.Name = fname.New(args[0])
	c.options.App.Bin = filepath.Join("bin", c.options.App.Name.String())

	if c.options.App.Name.String() == "." {
		c.options.App.Name = fname.New(filepath.Base(c.options.App.Root))
	} else {
		c.options.App.Root = filepath.Join(c.options.App.Root, c.options.App.Name.File().String())
	}

	c.options.App.PackageRoot(c.Module)
	if len(c.Module) == 0 {
		aa := meta.New(c.options.App.Root)
		c.options.App.PackageRoot(aa.PackagePkg)
	}

	if c.options.App.WithPop {
		if c.options.CI.DBType == "sqlite3" {
			c.options.App.WithSQLite = true
		}

		c.options.Pop = &pop.Options{
			Prefix:  c.options.App.Name.File().String(),
			Dialect: c.options.CI.DBType,
		}
	}

	run := genny.WetRunner(ctx)
	lg := logger.New(logger.DebugLevel)
	run.Logger = lg

	if c.DryRun {
		run = genny.DryRunner(ctx)
	}

	run.Root = c.options.App.Root
	if c.Force {
		os.RemoveAll(c.options.App.Root)
	}

	// initialize as API and then check if it's a web app
	// to then change.
	gg, err := api.New(&api.Options{
		Options: c.options,
	})

	if err != nil {
		return err
	}

	if !c.options.App.AsAPI {
		wo := &web.Options{
			Options: c.options,
		}

		if c.options.App.WithWebpack {
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

	if err := run.WithNew(gogen.Fmt(c.options.App.Root)); err != nil {
		return err
	}

	// setup VCS last
	if c.options.VCS != nil {
		// add the VCS generator
		if err := run.WithNew(vcs.New(c.options.VCS)); err != nil {
			return err
		}
	}

	if err := run.Run(); err != nil {
		return err
	}

	run.Logger.Infof("Congratulations! Your application, %s, has been successfully generated!", c.options.App.Name)
	run.Logger.Infof("You can find your new application at: %v", c.options.App.Root)
	run.Logger.Info("Please read the README.md file in your new application for next steps on running your application.")

	return nil
}
