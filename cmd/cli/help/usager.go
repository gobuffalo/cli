package help

// Usager is a plugin that wants to provide usage instructions
type Usager interface {
	// Usage return the string with the details on
	// how to use the plugin. p.e. buffalo help [command]
	Usage() string
}
