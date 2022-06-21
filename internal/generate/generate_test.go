package generate_test

import (
	"context"
	"flag"
	"fmt"

	"github.com/gobuffalo/cli/cmd/cli/clio"
	"github.com/gobuffalo/cli/internal/generate"
)

type tplugin string

func (tg tplugin) Name() string {
	return string(tg)
}

func (tg tplugin) Aliases() []string {
	return []string{
		// Adding some aliases
		string(tg)[:1],
		string(tg)[:2],
	}
}

var _ generate.Generator = testGenerator{tplugin("something"), &clio.IO{}}

type testGenerator struct {
	tplugin
	IO *clio.IO
}

func (tg testGenerator) HelpText() string {
	return fmt.Sprintf("%v test help text", string(tg.tplugin))
}

func (tg testGenerator) Generate(ctx context.Context, pwd string, arg []string) error {
	fmt.Fprintf(tg.IO.Stdout(), "generating: %v\n", string(tg.tplugin))

	return nil
}

type testFlagParserGenerator struct {
	testGenerator

	args []string
	err  error
}

func (tg *testFlagParserGenerator) ParseFlags(args []string) (*flag.FlagSet, error) {
	tg.args = args

	return flag.NewFlagSet("aa", flag.ContinueOnError), tg.err
}
