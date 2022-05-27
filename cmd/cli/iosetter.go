package cli

import "io"

// IOSetter is a plugin that allows its IO to be set
// the CLI will call this function with the IO.
type IOSetter interface {
	Plugin

	SetIO(stdin io.Reader, stdout, stderr io.Writer)
}
