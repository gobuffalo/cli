package test

import (
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/meta"
)

var (
	// if GO_BIN is set, use it, otherwise use "go"
	goBinary = envy.Get("GO_BIN", "go")

	// testValuedFlags are the flags on the go test command that
	// receive a value, we use this to determine if the last argument
	// is a flag or a package path
	testValuedFlags = strings.Join([]string{
		`-p`, `-asmflags`, `-buildmode`, `-compiler`, `-gccgoflags`,
		`-gcflags`, `-installsuffix`, `-mod`, `-modfile`, `-overlay`,
		`-pkgdir`, `-tags`, `-trimpath`, `-toolexec`, `-o`, `-exec`,
		`-bench`, `-benchtime`, `-blockprofile`, `-blockprofilerate`,
		`-count`, `-coverprofile`, `-cpu`, `-cpuprofile`, `-fuzz `,
		`-fuzzcachedir`, `-fuzzminimizetime`, `-fuzztime`, `-list`,
		`-memprofile`, `-memprofilerate`, `-mutexprofile`, `-mutexprofilefraction`,
		`-outputdir`, `-parallel`, `-run`, `-shuffle`, `-testlogfile`,
		`-timeout`, `-trace`,
	}, "|")
)

// buildCmd builds the test command to be passed to the go test command
// after cutting some of the arguments that are specific to buffalo and adding
// packages if missing.
func buildCmd(args []string) (*exec.Cmd, error) {
	app := meta.New(".")

	ccar := clean(args)
	cargs := append([]string{
		"test",
		// run sequentially
		"-p", "1",
		// add build tags
		"-tags", app.BuildTags("development").String(),
		//TODO: Should we merge it with passed tags?
	}, ccar...)

	// if no packages are specified, add the current directory
	lastIsFlag := len(ccar) > 0 && strings.HasPrefix(args[len(ccar)-1], "-")
	lastIsFlagValue := len(ccar) >= 2 && strings.Contains(testValuedFlags, ccar[len(ccar)-2])
	if len(ccar) == 0 || lastIsFlag || lastIsFlagValue {
		pkgs, err := findPackages()
		if err != nil {
			return nil, err
		}

		cargs = append(cargs, pkgs...)
	}

	// Add extra args (-testify.m) to the command
	cargs = append(cargs, extraArgs(args)...)

	cmd := exec.Command(goBinary, cargs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
}

// Removes flags that are only used by the Buffalo test command and should not be passed to
// the Go test command. Those flags are:
// * --force-migrations
// * -m
// * -testify.m
func clean(args []string) []string {
	var cleaned []string
	for ix, v := range args {
		if v == "--force-migrations" {
			continue
		}

		if v == "-m" || (ix > 0 && !strings.HasPrefix(v, "-") && args[ix-1] == "-m") {
			continue
		}

		if v == "-testify.m" || (ix > 0 && !strings.HasPrefix(v, "-") && args[ix-1] == "-testify.m") {
			continue
		}

		cleaned = append(cleaned, v)
	}

	return cleaned
}

// Adds extra flags to the go test command for the testify functionality if the
// -testify.m or -m flags are passed.
func extraArgs(args []string) []string {
	var extra []string
	for ix, v := range args {
		if ix >= len(args)-1 || strings.HasPrefix(args[ix+1], "-") {
			continue
		}

		if v == "-m" || v == "-testify.m" || v == "--run" {
			extra = append(extra, "-testify.m", args[ix+1])
		}
	}

	return extra
}
