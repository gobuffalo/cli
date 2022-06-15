package plugin

// Receiver allows a plugin to receive plugins.
// This interface is important for the extensibility of the CLI, as
// adding custom plugins will imply satisfying an interface that
// a PluginsReceiver is expecting so that the CLI can be extended.
type Receiver interface {
	// Name of the receiver.
	Name() string

	// Receive the plugins passed by the CLI, in here
	// a plugin can classify depending on then interfaces it
	// needs to satisfy.
	Receive(Plugins)
}
