package clio

import "io"

// IOSetter will be most of the commands for testing purposes.
type Setter interface {
	SetIO(stdout io.Writer, stderr io.Writer, stdin io.Reader)
}
