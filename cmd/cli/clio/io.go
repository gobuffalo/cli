package clio

import (
	"io"
	"os"
)

// Container represents the standard input, output, and error stream.
type Container struct {
	In  io.Reader // standard input
	Out io.Writer // standard output
	Err io.Writer // standard error
}

// Stdout returns IO.In.
// Defaults to os.Stdout.
func (oi Container) Stdout() io.Writer {
	if oi.Out == nil {
		return os.Stdout
	}

	return oi.Out
}

// Stderr returns IO.Err.
// Defaults to os.Stderr.
func (oi Container) Stderr() io.Writer {
	if oi.Err == nil {
		return os.Stderr
	}

	return oi.Err
}

// Stdin returns IO.In.
// Defaults to os.Stdin.
func (oi Container) Stdin() io.Reader {
	if oi.In == nil {
		return os.Stdin
	}

	return oi.In
}
