package cli

// FlagParser is a command that requires parsing flags,
// so the CLI needs to pass and parse args before calling Run.
type FlagParser interface {
	Plugin

	// ParseFlags and return the non-flag args.
	ParseFlags([]string) ([]string, error)
}
