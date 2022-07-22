package cli

import "github.com/gobuffalo/cli/cmd/cli/plugin"

// Adds received plugins to the list of plugins
// that the app has.
func (a *app) Add(pls ...plugin.Plugin) {
	a.plugins = append(a.plugins, pls...)
}

// Removes plugins from the list of plugins that
// match the given names.
func (a *app) Remove(names ...string) {
	result := plugin.Plugins{}

	for _, v := range a.plugins {
		var found bool
		for _, x := range names {
			found = found || x == v.Name()
		}

		if !found {
			result = append(result, v)
		}
	}

	a.plugins = result
}

// Clears the list of plugins, keeps `help`` and
// `plugins` commands.
func (a *app) Clear() {
	a.plugins = basePlugins
}
