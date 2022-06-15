package cli

// The Plugin interface is used to identify plugins
// that can be loaded into the CLI. Plugins get specific
// depending on its usage. p.e. Command or BeforeTester.
type Plugin interface {
	Name() string
}
