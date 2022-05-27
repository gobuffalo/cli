package cli

import (
	"context"
)

// Command is a plugin that can Run
type Command interface {
	Plugin

	Run(ctx context.Context, pwd string, args []string) error
}
