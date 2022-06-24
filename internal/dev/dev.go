package dev

import (
	"context"
)

func contextAwareRun(ctx context.Context, f func() error) error {
	var out = make(chan error)

	go func() {
		out <- f()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-out:
		return err
	}
}
