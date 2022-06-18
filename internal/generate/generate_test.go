package generate_test

import (
	"context"
	"fmt"
)

type testGenerator string

func (tg testGenerator) Name() string {
	return string(tg)
}

func (tg testGenerator) Aliases() []string {
	return []string{
		// Adding some aliases
		string(tg)[:1],
		string(tg)[:2],
	}
}

func (tg testGenerator) HelpText() string {
	return fmt.Sprintf("%v test help text", string(tg))
}

func (tg testGenerator) Generate(context.Context, string, []string) error {
	return nil
}
