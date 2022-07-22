package cli

import "github.com/gobuffalo/cli/cmd/cli/plugin"

// Adds received plugins to the list of plugins
// that the app has.
func (app *App) Add(pls ...plugin.Plugin) {
	app.plugins = append(app.plugins, pls...)
}

// Removes plugins from the list of plugins that
// match the given names.
func (app *App) Remove(names ...string) {
	result := plugin.Plugins{}

	for _, v := range app.plugins {
		var found bool
		for _, x := range names {
			found = found || x == v.Name()
		}

		if !found {
			result = append(result, v)
		}
	}

	app.plugins = result
}

// Clears the list of plugins, keeps `help`` and
// `plugins` commands.
func (app *App) Clear() {
	app.plugins = basePlugins
}