package cli

// DefaultApp is an instance of the CLI application
// loaded with `default` plugins. The `NewApp` function
// could be used to create a custom instance of the CLI
// with custom plugins.
var DefaultApp = NewApp(defaultPlugins...)
