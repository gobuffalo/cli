package dev

import (
	"context"
	"os"

	rg "github.com/gobuffalo/cli/internal/genny/refresh"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/refresh/refresh"
)

var StartServer = &startServer{}

type startServer struct {
	debug bool
}

func (ss startServer) Name() string {
	return "dev/start-server"
}

func (ss startServer) HelpText() string {
	return "Starts the development server with Go."
}

func (ss *startServer) EnableDebug() {
	ss.debug = true
}

func (ss startServer) RunDevelopment(ctx context.Context, app meta.App, args []string) error {
	cfgFile := "./.buffalo.dev.yml"
	if _, err := os.Stat(cfgFile); err != nil {
		run := genny.WetRunner(ctx)
		err = run.WithNew(rg.New(&rg.Options{App: app}))
		if err != nil {
			return err
		}

		if err := run.Run(); err != nil {
			return err
		}
	}

	c := &refresh.Configuration{}
	if err := c.Load(cfgFile); err != nil {
		return err
	}

	c.Debug = ss.debug

	bt := app.BuildTags("development")
	for _, v := range bt {
		c.BuildFlags = append(c.BuildFlags, "-tags", v)
	}

	r := refresh.NewWithContext(c, ctx)
	r.CommandFlags = args

	return contextAwareRun(ctx, r.Start)
}
