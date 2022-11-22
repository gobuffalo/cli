package test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/pop/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func runE(c *cobra.Command, args []string) error {
	os.Setenv("GO_ENV", "test")
	if _, err := os.Stat("database.yml"); err == nil {
		// there's a database
		test, err := pop.Connect("test")
		if err != nil {
			return err
		}

		// drop the test db:
		if err := test.Dialect.DropDB(); err != nil {
			// not an error, since the database will be created in the next step anyway
			logrus.Info("no test database to drop")
		}

		// create the test db:
		err = test.Dialect.CreateDB()
		if err != nil {
			return err
		}

		// Read and remove --force-migrations flag from args:
		forceMigrations := strings.Contains(strings.Join(args, ""), "--force-migrations")
		args = cutArg("--force-migrations", args)
		if forceMigrations {
			fm, err := pop.NewFileMigrator("./migrations", test)
			if err != nil {
				return err
			}

			if err := fm.Up(); err != nil {
				return err
			}

			return testRunner(args)
		}

		if schema := findSchema(); schema != nil {
			err = test.Dialect.LoadSchema(schema)
			if err != nil {
				return err
			}
		}
	}
	return testRunner(args)
}

func findSchema() io.Reader {
	if f, err := os.Open(filepath.Join("migrations", "schema.sql")); err == nil {
		return f
	}
	if dev, err := pop.Connect("development"); err == nil {
		schema := &bytes.Buffer{}
		if err = dev.Dialect.DumpSchema(schema); err == nil {
			return schema
		}
	}

	if test, err := pop.Connect("test"); err == nil {
		fm, err := pop.NewFileMigrator("./migrations", test)
		if err != nil {
			return nil
		}

		if err := fm.Up(); err == nil {
			return nil
		}
	}
	return nil
}

func testRunner(args []string) error {
	var mFlag bool
	var query string

	commandArgs := []string{}
	packageArgs := []string{}

	var lastArg string
	for index, arg := range args {
		switch arg {
		case "-run", "-m":
			query = args[index+1]
			mFlag = true
		case "-v", "-timeout":
			commandArgs = append(commandArgs, arg)
		default:
			if lastArg == "-timeout" {
				commandArgs = append(commandArgs, arg)
			} else if lastArg != "-run" && lastArg != "-m" {
				packageArgs = append(packageArgs, arg)
			}
		}

		lastArg = arg
	}

	cmd := newTestCmd(commandArgs)
	if mFlag {
		return mFlagRunner{
			query: query,
			args:  commandArgs,
			pargs: packageArgs,
		}.Run()
	}

	pkgs, err := findTestPackages(packageArgs)
	if err != nil {
		return err
	}

	cmd.Args = append(cmd.Args, pkgs...)
	logrus.Info(strings.Join(cmd.Args, " "))
	return cmd.Run()
}

type mFlagRunner struct {
	query string
	args  []string
	pargs []string
}

func (m mFlagRunner) Run() error {
	app := meta.New(".")
	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)

	pkgs, err := findTestPackages(m.pargs)
	if err != nil {
		return err
	}

	var errs bool
	for _, p := range pkgs {
		os.Chdir(pwd)

		if p == app.PackagePkg {
			continue
		}

		cmd := newTestCmd(m.args)

		p = strings.TrimPrefix(p, app.PackagePkg+string(filepath.Separator))
		os.Chdir(p)

		if hasTestify(cmd.Args) {
			cmd.Args = append(cmd.Args, "-testify.m", m.query)
		} else {
			cmd.Args = append(cmd.Args, "-run", m.query)
		}

		logrus.Info(strings.Join(cmd.Args, " "))

		if err := cmd.Run(); err != nil {
			errs = true
		}
	}
	if errs {
		return fmt.Errorf("errors running tests")
	}
	return nil
}

func hasTestify(args []string) bool {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Args = append(cmd.Args, "-unknownflag")
	b, _ := cmd.Output()
	return bytes.Contains(b, []byte("-testify.m"))
}

func newTestCmd(args []string) *exec.Cmd {
	cargs := []string{"test", "-p", "1"}
	app := meta.New(".")
	cargs = append(cargs, "-tags", app.BuildTags("development").String())
	cargs = append(cargs, args...)
	cmd := exec.Command(envy.Get("GO_BIN", "go"), cargs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func cutArg(arg string, args []string) []string {
	for i, v := range args {
		if v == arg {
			return append(args[:i], args[i+1:]...)
		}
	}

	return args
}
