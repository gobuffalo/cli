package plugdeps

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/meta"
)

// ErrMissingConfig is if config/buffalo-plugins.toml file is not found. Use plugdeps#On(app) to test if plugdeps are being used
var ErrMissingConfig = fmt.Errorf("could not find a buffalo-plugins config file at %s", ConfigPath(meta.New(".")))

// List all of the plugins the application depeneds on. Will return ErrMissingConfig
// if the app is not using config/buffalo-plugins.toml to manage their plugins.
// Use plugdeps#On(app) to test if plugdeps are being used.
func List(app meta.App) (*Plugins, error) {
	plugs := New()

	lp, err := listLocal(app)
	if err != nil {
		return plugs, err
	}
	plugs.Add(lp.List()...)

	if !On(app) {
		return plugs, ErrMissingConfig
	}

	p := ConfigPath(app)
	tf, err := os.Open(p)
	if err != nil {
		return plugs, err
	}
	if err := plugs.Decode(tf); err != nil {
		return plugs, err
	}

	return plugs, nil
}

func listLocal(app meta.App) (*Plugins, error) {
	plugs := New()
	proot := filepath.Join(app.Root, "plugins")
	if _, err := os.Stat(proot); err != nil {
		return plugs, nil
	}

	err := fs.WalkDir(os.DirFS(proot), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		if strings.HasPrefix(base, "buffalo-") {
			plugs.Add(Plugin{
				Binary: base,
				Local:  "." + strings.TrimPrefix(path, app.Root),
			})
		}
		return nil
	})
	return plugs, err
}

// ConfigPath returns the path to the config/buffalo-plugins.toml file
// relative to the app
func ConfigPath(app meta.App) string {
	return filepath.Join(app.Root, "config", "buffalo-plugins.toml")
}

// On checks for the existence of config/buffalo-plugins.toml if this
// file exists its contents will be used to list plugins. If the file is not
// found, then the BUFFALO_PLUGIN_PATH and ./plugins folders are consulted.
func On(app meta.App) bool {
	_, err := os.Stat(ConfigPath(app))
	return err == nil
}
