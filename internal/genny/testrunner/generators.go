package testrunner

import (
	"github.com/gobuffalo/cli/internal/genny/newapp/api"
	"github.com/gobuffalo/cli/internal/genny/newapp/web"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gentest"
)

func WebApp() (*genny.Runner, error) {
	gg, err := web.New(&web.Options{})
	if err != nil {
		return nil, err
	}
	return newApp(gg)
}

func ApiApp() (*genny.Runner, error) {
	gg, err := api.New(&api.Options{})
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
