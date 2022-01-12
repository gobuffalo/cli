package new

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	pop "github.com/gobuffalo/buffalo-pop/v3/genny/newapp"
	"github.com/gobuffalo/cli/internal/genny/ci"
	"github.com/gobuffalo/cli/internal/genny/docker"
	"github.com/gobuffalo/cli/internal/genny/newapp/core"
	"github.com/gobuffalo/cli/internal/genny/refresh"
	"github.com/gobuffalo/cli/internal/genny/vcs"
	fname "github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/meta"
	"github.com/spf13/viper"
)

var configError error

type options struct {
	Options *core.Options
	Module  string
	Force   bool
	Verbose bool
	DryRun  bool
}

func parseNewOptions(args []string) (options, error) {
	nopts := options{
		Force:   viper.GetBool("force"),
		Verbose: viper.GetBool("verbose"),
		DryRun:  viper.GetBool("dry-run"),
		Module:  viper.GetString("module"),
	}

	if len(args) == 0 {
		return nopts, fmt.Errorf("you must enter a name for your new application")
	}
	if configError != nil {
		return nopts, configError
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nopts, err
	}
	app := meta.New(pwd)
	app.WithGrifts = true
	app.Name = fname.New(args[0])
	app.Bin = filepath.Join("bin", app.Name.String())

	if app.Name.String() == "." {
		app.Name = fname.New(filepath.Base(app.Root))
	} else {
		app.Root = filepath.Join(app.Root, app.Name.File().String())
	}

	if len(nopts.Module) == 0 {
		aa := meta.New(app.Root)
		app.PackageRoot(aa.PackagePkg)
	} else {
		app.PackageRoot(nopts.Module)
	}

	app.AsAPI = viper.GetBool("api")
	app.VCS = viper.GetString("vcs")

	app.WithPop = !viper.GetBool("skip-pop")
	app.WithWebpack = !viper.GetBool("skip-webpack")
	app.WithYarn = !viper.GetBool("skip-yarn")
	app.WithNodeJs = app.WithWebpack
	app.AsWeb = !app.AsAPI

	if app.AsAPI {
		app.WithWebpack = false
		app.WithYarn = false
		app.WithNodeJs = false
	}

	opts := &core.Options{}

	if x := viper.GetBool("skip-docker"); !x {
		opts.Docker = &docker.Options{}
	}

	app.WithDocker = !viper.GetBool("skip-docker")

	if x := viper.GetString("ci-provider"); len(x) > 0 && x != "none" {
		opts.CI = &ci.Options{
			Provider: x,
			DBType:   viper.GetString("db-type"),
		}
	}

	if len(app.VCS) > 0 && app.VCS != "none" {
		opts.VCS = &vcs.Options{
			Provider: app.VCS,
		}
	}

	if app.WithPop {
		d := viper.GetString("db-type")
		if d == "sqlite3" {
			app.WithSQLite = true
		}

		opts.Pop = &pop.Options{
			Prefix:  app.Name.File().String(),
			Dialect: d,
		}
	}

	opts.Refresh = &refresh.Options{}

	opts.App = app
	nopts.Options = opts

	return nopts, nil
}

func initConfig(skipConfig *bool, cfgFile *string) func() {
	return func() {
		if *skipConfig {
			return
		}

		var err error
		if *cfgFile != "" { // enable ability to specify config file via flag
			viper.SetConfigFile(*cfgFile)
			// Will error only if the --config flag is used
			if err = viper.ReadInConfig(); err != nil {
				configError = err
			}
		} else {
			viper.SetConfigName(".buffalo") // name of config file (without extension)
			viper.AddConfigPath("$HOME")    // adding home directory as first search path
			viper.AutomaticEnv()            // read in environment variables that match
			viper.ReadInConfig()
		}
	}
}

func currentUser() (string, error) {
	_, err := exec.LookPath("git")
	if err == nil {
		b, err := exec.Command("git", "config", "github.user").Output()
		if err == nil {
			return string(b), nil
		}
	}

	u, err := user.Current()
	if err != nil {
		return "", err
	}

	username := u.Username
	if t := strings.Split(username, `\`); len(t) > 0 {
		username = t[len(t)-1]
	}

	return username, nil
}
