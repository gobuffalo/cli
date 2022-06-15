package cli

// Aliaser is a plugin that defines a list of aliases
// to be identified.
type Aliaser interface {
	Plugin

	Aliases() []string
}
