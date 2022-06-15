package plugin

// Aliaser is a plugin that defines a list of aliases
// to be identified.
type Aliaser interface {
	// Name of the plugin
	Name() string

	// Aliases of the plugin returned in a list
	// of strings.
	Aliases() []string
}
