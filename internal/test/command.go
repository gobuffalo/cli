package test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/cli/internal/tools/pop"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/meta"
)

// The test command instance with default before and
// after test plugins.
var Command = &command{
	before: []BeforeTester{
		TestEnvironment,
		&pop.SetupTestDatabase{},
	},

	testers: []Tester{
		defaultTester("Buffalo"),
	},
}

type command struct {
	before  []BeforeTester
	testers []Tester
}

func (c command) Name() string {
	return "test"
}

func (c command) HelpText() string {
	// TODO: List before and after plugins
	return "Runs application tests by invoking before and after test plugins."
}

func (c *command) Main(ctx context.Context, pwd string, args []string) error {
	// Iterate over the BeforeTesters and run each of them
	// in case of an error halt the testing process by returning the error.
	for _, v := range c.before {
		err := v.BeforeTest(ctx, pwd, args)
		if err == nil {
			continue
		}

		return fmt.Errorf("error running `%s` before test: %w", v.Name(), err)
	}

	// Running the testers, Go, JS and any others could go here.
	for _, t := range c.testers {
		err := t.Test(ctx, pwd, args)
		if err != nil {
			return fmt.Errorf("error running `%s` test: %w", t.Name(), err)
		}
	}

	return nil
}

func hasTestify(p string) bool {
	cmd := exec.Command("go", "test", "-thisflagdoesntexist")
	b, _ := cmd.Output()

	return bytes.Contains(b, []byte("-testify.m"))
}

func testPackages(givenArgs []string) ([]string, error) {
	// If there are args, then assume these are the packages to test.
	//
	// Instead of always returning all packages from 'go list ./...', just
	// return the given packages in this case
	if len(givenArgs) > 0 {
		return givenArgs, nil
	}

	args := []string{}
	out, err := exec.Command(envy.Get("GO_BIN", "go"), "list", "./...").Output()
	if err != nil {
		return args, err
	}
	pkgs := bytes.Split(bytes.TrimSpace(out), []byte("\n"))
	for _, p := range pkgs {
		if !strings.Contains(string(p), "/vendor/") {
			args = append(args, string(p))
		}
	}
	return args, nil
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
