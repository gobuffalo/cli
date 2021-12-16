package fix

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/gobuffalo/cli/internal/genny/plugins/install"

	cmdPlugins "github.com/gobuffalo/cli/internal/cmd/plugins"
	"github.com/gobuffalo/cli/internal/plugins"
	"github.com/gobuffalo/cli/internal/plugins/plugdeps"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/meta"
)

var current_pop = "github.com/gobuffalo/buffalo-pop/v3"

//Plugins fixes the plugin configuration of the project by
//manipulating the plugins .toml file.
type Plugins struct{}

//CleanCache cleans the plugins cache folder by removing it
func (pf Plugins) CleanCache(r *Runner) error {
	fmt.Println("~~~ Cleaning plugins cache ~~~")
	os.RemoveAll(plugins.CachePath)
	return nil
}

//Reinstall installs latest versions of the plugins
func (pf Plugins) Reinstall(r *Runner) error {
	plugs, err := plugdeps.List(r.App)
	if err != nil && !errors.Is(err, plugdeps.ErrMissingConfig) {
		return err
	}

	// TODO: generalize with meta/v2
	if r.App.WithPop {
		plugs.Add(plugdeps.NewPlugin(current_pop))
	}
	run := genny.WetRunner(context.Background())
	gg, err := install.New(&install.Options{
		App:     r.App,
		Plugins: plugs.List(),
	})
	if err != nil {
		return err
	}

	run.WithGroup(gg)

	fmt.Println("~~~ Reinstalling plugins ~~~")
	return run.Run()
}

//RemoveOld removes old and deprecated plugins
func (pf Plugins) RemoveOld(r *Runner) error {
	fmt.Println("~~~ Removing old plugins ~~~")

	run := genny.WetRunner(context.Background())
	app := meta.New(".")
	plugs, err := plugdeps.List(app)
	if err != nil && !errors.Is(err, plugdeps.ErrMissingConfig) {
		return err
	}

	plugs.Remove(plugdeps.Plugin{
		Binary: "buffalo-pop",
	})

	fmt.Println("~~~ Removing previous version of buffalo-pop plugin ~~~")

	run.WithRun(cmdPlugins.NewEncodePluginsRunner(app, plugs))

	return run.Run()
}
