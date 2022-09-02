package help

// HelpTexter is a command that provides a help text method
// which could be used on the app Usage method.
type HelpTexter interface {
	HelpText() string
}
