package dev

import (
	"context"

	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

type DevelopmentSetupper interface {
	plugin.Plugin

	SetupDevelopment(context.Context, string, []string) error
}
