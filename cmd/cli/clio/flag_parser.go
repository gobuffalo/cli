package clio

import flag "github.com/spf13/pflag"

// FlagParser interface allows commands to parse flags according to their needs.
type FlagParser interface {
	// ParseFlags receives the args and returns a pointer to the parsed FlagSet.
	ParseFlags(args []string) (*flag.FlagSet, error)
}
