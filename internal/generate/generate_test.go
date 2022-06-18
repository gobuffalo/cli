package generate_test

import (
	"context"
	"fmt"
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

type testGenerator struct {
	tplugin
}

func (tg testGenerator) HelpText() string {
	return fmt.Sprintf("%v test help text", string(tg.tplugin))
}

func (tg testGenerator) Generate(context.Context, string, []string) error {
	return nil
}
