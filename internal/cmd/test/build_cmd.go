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
// TODO: -run -testify.m -testify.m
func buildCmd(args []string) *exec.Cmd {
	app := meta.New(".")

	if len(args) > 0 {
		// Cleanup --force-migrations flag
		args = strings.Split(strings.Replace(strings.Join(args, " "), "--force-migrations ", "", -1), " ")
	}

	cargs := append([]string{
		"test",
		// run sequentially
		"-p", "1",
		// add build tags
		"-tags", app.BuildTags("development").String(),
		//TODO: Should we merge it with passed tags?
	}, args...)

	// if no packages are specified, add the current directory
	lastIsFlag := len(args) > 0 && strings.HasPrefix(args[len(args)-1], "-")
	lastIsFlagValue := len(args) >= 2 && strings.Contains(testValuedFlags, args[len(args)-2])
	if len(args) == 0 || lastIsFlag || lastIsFlagValue {
		// TODO: Should we run one by one the packages?
		cargs = append(cargs, "./...")
	}

	cmd := exec.Command(goBinary, cargs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
