package new

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"

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

	force       bool
	verbose     bool
	dryRun      bool
	skipPop     bool
	skipWebpack bool
	skipYarn    bool
	skipDocker  bool
	api         bool

	module     string
	vcs        string
	dbType     string
	ciProvider string
}

func (c command) Name() string {
	return "new"
}

func (c command) HelpText() string {
	return "Creates a new Buffalo application"
}

func (c command) Usage() string {
	return "buffalo new [flags] [application-name]"
}

func (c *command) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet("new", flag.ContinueOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(io.Discard)
	}

	c.flagSet.BoolVar(&c.api, "api", false, "skip all front-end code and configure for an API server")
	c.flagSet.BoolVar(&c.skipPop, "skip-pop", false, "skip all back-end code and configure for a web server")
	c.flagSet.BoolVar(&c.skipWebpack, "skip-webpack", false, "skip all front-end code and configure for a web server")
	c.flagSet.BoolVar(&c.skipYarn, "skip-yarn", false, "skip all front-end code and configure for a web server")
	c.flagSet.BoolVar(&c.skipDocker, "skip-docker", false, "skip all front-end code and configure for a web server")
	c.flagSet.BoolVarP(&c.force, "force", "f", false, "delete and remake if the app already exists")
	c.flagSet.BoolVarP(&c.dryRun, "dry-run", "d", false, "dry run")
	c.flagSet.BoolVarP(&c.verbose, "verbose", "v", false, "verbosely print out the go get commands")

	c.flagSet.StringVar(&c.module, "module", "", "module to use for the application")
	c.flagSet.StringVar(&c.vcs, "vcs", "git", "specify the Version control system you would like to use [none, git, bzr]")
	c.flagSet.StringVar(&c.dbType, "db", "postgres", fmt.Sprintf("specify the type of database you want to use [%s]", strings.Join(plib.AvailableDialects, ", ")))
	c.flagSet.StringVar(&c.ciProvider, "ci-provider", "travis", "specify the CI provider you want to use [none, travis, gitlab-ci, circleci]")

	// if len(args) >= 1 && !strings.HasPrefix(args[0], "-") {
	// 	fmt.Println("Usage: " + c.Usage())

	// 	return c.flagSet, fmt.Errorf("error: flags must go before the application name")
	// }

	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}

func (c *command) Main(ctx context.Context, pwd string, args []string) error {
	fmt.Println(args)

	args = c.flagSet.Args()

	if len(args) == 0 {
		return fmt.Errorf("you must enter a name for your new application")
	}

	c.options.App = c.buildApp(pwd, args[0])
	c.setOptions()

	run := genny.WetRunner(ctx)
	if c.dryRun {
		run = genny.DryRunner(ctx)
	}
	// Setting debug logger if verbose is set
	if c.verbose {
		run.Logger = logger.New(logger.DebugLevel)
	}

	// Remove existing folder if Force is passed.
	run.Root = c.options.App.Root
	if c.force {
		// TODO: this needs to considerate the -dry-run flag.
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

	fmt.Printf("\nCongratulations! Your application, %s, has been successfully generated!\n", c.options.App.Name)
	fmt.Printf("You can find your new application at: %v\n", c.options.App.Root)
	fmt.Printf("Please read the README.md file in your new application for next steps on running your application.\n")

	return nil
}

func (c command) buildApp(wd string, name string) meta.App {
	app := meta.New(wd)

	app.Name = fname.New(name)
	app.Bin = filepath.Join("bin", app.Name.String())

	app.WithDocker = !c.skipDocker
	app.WithWebpack = !c.skipWebpack
	app.WithYarn = !c.skipYarn
	app.WithPop = !c.skipPop
	app.AsAPI = c.api

	app.WithGrifts = true
	app.WithNodeJs = c.options.App.WithWebpack
	app.AsWeb = !c.options.App.AsAPI

	if app.AsAPI {
		app.WithWebpack = false
		app.WithYarn = false
		app.WithNodeJs = false
	}

	if app.Name.String() == "." {
		app.Name = fname.New(filepath.Base(app.Root))
	} else {
		app.Root = filepath.Join(app.Root, app.Name.File().String())
	}

	app.PackageRoot(c.module)
	if len(c.module) == 0 {
		aa := meta.New(app.Root)
		app.PackageRoot(aa.PackagePkg)
	}

	if app.WithPop && c.dbType == "sqlite3" {
		app.WithSQLite = true
	}

	c.options.VCS = &vcs.Options{
		Provider: c.vcs,
	}

	return app
}

func (c *command) setOptions() {
	c.options.Refresh = &refresh.Options{}

	if c.options.App.WithDocker {
		c.options.Docker = &docker.Options{}
	}

	if c.ciProvider != "none" {
		c.options.CI = &ci.Options{
			Provider: c.ciProvider,
			DBType:   c.dbType,
		}
	}

	if c.options.App.WithPop {
		c.options.Pop = &pop.Options{
			Prefix:  c.options.App.Name.File().String(),
			Dialect: c.dbType,
		}
	}
}
