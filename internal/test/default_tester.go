package test

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gobuffalo/envy"
)

type defaultTester string

func (dt defaultTester) Name() string {
	return "buffalo/tester"
}

func (dt defaultTester) HelpText() string {
	return "buffalo/tester"
}

func (dt defaultTester) Test(ctx context.Context, pwd string, args []string) error {
	// Start with the actual testing process.
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

	// If there are args, then assume these are the packages to test.
	//
	// Instead of always returning all packages from 'go list ./...', just
	// return the given packages in this case
	pkgs := []string{}
	if len(packageArgs) == 0 {
		out, err := exec.Command(envy.Get("GO_BIN", "go"), "list", "./...").Output()
		if err != nil {
			return err
		}

		pkgs = strings.Split(strings.TrimSpace(string(out)), "\n")
	}

	cmd.Args = append(cmd.Args, pkgs...)
	fmt.Printf("\nRunning: %s \n\n", strings.Join(cmd.Args, " "))

	return cmd.Run()
}
