package generate_test

import (
	"strings"
	"testing"

	"github.com/gobuffalo/cli/cmd/cli/plugin"
	"github.com/gobuffalo/cli/internal/generate"
)

func TestCommandLongHelpText(t *testing.T) {
	t.Run("no generators", func(t *testing.T) {
		var lt = generate.Command.LongHelpText()
		if strings.Contains(lt, "Registered Generators") {
			t.Fatalf("expected not to find `Registered Generators` in the long help text")
		}

		if !strings.Contains(lt, "No generators registered") {
			t.Fatalf("expected to find `No generators registered` in the long help text")
		}
	})

	t.Run("generator added", func(t *testing.T) {
		generate.Command.Receive(plugin.Plugins{
			testGenerator("something"),
			testGenerator("other"),
		})

		exp := []string{
			"Registered Generators",
			"something",
			"other",
			"something test help text",
			"other test help text",
			"s, so",
			"o, ot",
		}

		var lt = generate.Command.LongHelpText()
		for _, e := range exp {
			if strings.Contains(lt, e) {
				continue
			}

			t.Fatalf("expected to find `%v` in the long help text", e)
		}

	})
}

func TestCommandHelp(t *testing.T) {

}
