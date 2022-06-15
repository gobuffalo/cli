package help

// LongHelpTexter is a command that provides a help text method
// which could be used on the app Usage method.
type LongHelpTexter interface {
	LongHelpText() string
}
