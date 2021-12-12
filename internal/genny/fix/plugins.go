package fix

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gobuffalo/cli/internal/genny/plugins/install"

	cmdPlugins "github.com/gobuffalo/cli/internal/cmd/plugins"
	"github.com/gobuffalo/cli/internal/plugins"
	"github.com/gobuffalo/cli/internal/plugins/plugdeps"
	"github.com/gobuffalo/genny/v2"
)

var oldPlugins = []string{
	"github.com/gobuffalo/buffalo-pop",
	"github.com/gobuffalo/buffalo-pop/v2",
}

// CleanPluginCache cleans the plugins cache folder by removing it
func CleanPluginCache(r *genny.Runner) error {
	fmt.Println("~~~ Cleaning plugins cache ~~~")
	os.RemoveAll(plugins.CachePath)
	return nil
}

// ReinstallPlugins installs latest versions of the plugins
func ReinstallPlugins(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		plugs, err := plugdeps.List(opts.App)
		if err != nil && !errors.Is(err, plugdeps.ErrMissingConfig) {
			return err
		}

		fmt.Println("~~~ Reinstalling plugins ~~~")

		gg, err := install.New(&install.Options{
			App:     opts.App,
			Plugins: plugs.List(),
		})
		if err != nil {
			return err
		}

		r.WithGroup(gg)
		return nil
	}
}

// RemoveOldPlugins removes old and deprecated plugins
func RemoveOldPlugins(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		fmt.Println("~~~ Removing old plugins ~~~")

		plugs, err := plugdeps.List(opts.App)
		if err != nil && !errors.Is(err, plugdeps.ErrMissingConfig) {
			return err
		}

		for _, p := range oldPlugins {
			a := strings.TrimSpace(p)
			bin := path.Base(a)
			plugs.Remove(plugdeps.Plugin{
				Binary: bin,
				GoGet:  a,
			})

			fmt.Println("~~~ Removing", p, "plugin ~~~")
			r.WithRun(cmdPlugins.NewEncodePluginsRunner(opts.App, plugs))
		}
		return nil
	}
}
