package generate

import "context"

// Generator is a plugin that generates code block
type Generator interface {
	// Name of the generator, useful for the
	// help of the generate command.
	Name() string

	// HelpText is a short description of the
	// generator.
	HelpText() string

	// Generate the desired code block or return
	// and error if something goes wrong.
	Generate(context.Context, string, []string) error
}
