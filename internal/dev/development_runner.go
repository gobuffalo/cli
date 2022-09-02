package dev

import (
	"context"

	"github.com/gobuffalo/cli/cmd/cli/plugin"
	"github.com/gobuffalo/meta"
)

type DevelopmentRunner interface {
	plugin.Plugin

	EnableDebug()
	RunDevelopment(context.Context, meta.App, []string) error
}
