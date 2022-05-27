package cli

import (
	"io"
	"os"
)

// IO contains the IN, OUT and ERR streams for the CLI.
type IO struct {
	In  io.Reader // standard input
	Out io.Writer // standard output
	Err io.Writer // standard error
}

// Stdout returns IO.In.
// Defaults to os.Stdout.
func (oi IO) Stdout() io.Writer {
	if oi.Out == nil {
		return os.Stdout
	}

	return oi.Out
}

// Stderr returns IO.Err.
// Defaults to os.Stderr.
func (oi IO) Stderr() io.Writer {
	if oi.Err == nil {
		return os.Stderr
	}

	return oi.Err
}

// Stdin returns IO.In.
// Defaults to os.Stdin.
func (oi IO) Stdin() io.Reader {
	if oi.In == nil {
		return os.Stdin
	}

	return oi.In
}
