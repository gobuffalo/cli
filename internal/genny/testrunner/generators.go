package testrunner

import (
	"github.com/gobuffalo/cli/internal/genny/newapp/api"
	"github.com/gobuffalo/cli/internal/genny/newapp/web"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gentest"
)

func WebApp(opts *web.Options) (*genny.Runner, error) {
	gg, err := web.New(opts)
	if err != nil {
		return nil, err
	}
	return newApp(gg)
}

func ApiApp(opts *api.Options) (*genny.Runner, error) {
	gg, err := api.New(opts)
	if err != nil {
		return nil, err
	}
	return newApp(gg)
}

func newApp(gg *genny.Group) (*genny.Runner, error) {
	run := gentest.NewRunner()
	run.WithGroup(gg)
	if err := run.Run(); err != nil {
		return nil, err
	}

	runner := gentest.NewRunner()
	for _, f := range run.Results().Files {
		runner.Disk.Add(f)
	}
	return runner, nil
}
