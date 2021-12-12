package fix

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/gobuffalo/genny/v2"
)

func installTools(opts *Options) ([]string, error) {
	fmt.Println("~~~ Installing required tools ~~~")
	run := genny.WetRunner(context.Background())
	g := genny.New()
	if opts.App.WithPop {
		g.Command(exec.Command("go", "install", "github.com/gobuffalo/buffalo-pop/v3@latest"))
	}
	if err := run.With(g); err != nil {
		return nil, err
	}
	return nil, run.Run()
}
