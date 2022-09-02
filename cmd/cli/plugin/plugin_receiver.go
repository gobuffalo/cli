package plugin

// Receiver allows a plugin to receive plugins.
// This interface is important for the extensibility of the CLI, as
// adding custom plugins will imply satisfying an interface that
// a PluginsReceiver is expecting so that the CLI can be extended.
type Receiver interface {
	Plugin

	// Receive the plugins passed by the CLI, in here
	// a plugin can classify depending on then interfaces it
	// needs to satisfy. This function may be called multiple times
	// And is up to the receiver to handle the multiple set of plugins.
	Receive(Plugins)
}
