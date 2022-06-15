package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/meta"
)

type mFlagRunner struct {
	query string
	args  []string
	pargs []string
}

func (m mFlagRunner) Run() error {
	app := meta.New(".")

	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)

	pkgs, err := testPackages(m.pargs)
	if err != nil {
		return err
	}

	var errs bool
	for _, p := range pkgs {
		os.Chdir(pwd)

		if p == app.PackagePkg {
			continue
		}

		p = strings.TrimPrefix(p, app.PackagePkg+string(filepath.Separator))
		os.Chdir(p)

		cmd := newTestCmd(m.args)
		if hasTestify(p) {
			cmd.Args = append(cmd.Args, "-testify.m", m.query)
		} else {
			cmd.Args = append(cmd.Args, "-run", m.query)
		}

		if err := cmd.Run(); err != nil {
			errs = true
		}
	}
	if errs {
		return fmt.Errorf("errors running tests")
	}
	return nil
}
