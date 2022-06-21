package generate_test

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/gobuffalo/cli/cmd/cli/clio"
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
			testGenerator{
				tplugin: tplugin("something"),
			},
			testGenerator{
				tplugin: tplugin("other"),
			},
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

func TestReceive(t *testing.T) {
	t.Run("no generators", func(t *testing.T) {
		var g = generate.Command
		g.Receive(plugin.Plugins{
			tplugin("ax"),
			tplugin("bw"),
		})

		ht := g.LongHelpText()
		t.Log(ht)
		if !strings.Contains(ht, "No generators registered") {
			t.Fatalf("expected to find `No generators registered` in the long help text")
		}
	})

	t.Run("generators passed", func(t *testing.T) {
		var g = generate.Command
		g.Receive(plugin.Plugins{
			tplugin("ax"),
			tplugin("bw"),

			testGenerator{
				tplugin: tplugin("mw"),
			},
		})

		ht := g.LongHelpText()
		t.Log(ht)

		if strings.Contains(ht, "No generators registered") {
			t.Fatalf("expected not to find `No generators registered` in the long help text")
		}

		exp := []string{
			"mw",
			"m, mw",
			"Registered Generators",
		}

		for _, e := range exp {
			if !strings.Contains(ht, e) {
				t.Fatalf("expected to find `%v` in the long help text", e)
			}
		}
	})
}

func TestGenerateParseFlags(t *testing.T) {
	var g = generate.Command
	t.Run("calls generators with args", func(t *testing.T) {
		parser := &testFlagParserGenerator{
			testGenerator: testGenerator{
				tplugin: tplugin("mw"),
			},
		}

		g.Receive(plugin.Plugins{
			parser,
		})

		g.ParseFlags([]string{"tg", "foo", "bar"})

		if len(parser.args) != 2 {
			t.Fatalf("generator should only receive 2 args")
		}

		if parser.args[0] != "foo" {
			t.Fatalf("expected generator to receive the flags")
		}

	})

	t.Run("doesn't care about parsing errors", func(t *testing.T) {
		sec := &testFlagParserGenerator{
			testGenerator: testGenerator{tplugin: tplugin("xw")},
		}

		erp := &testFlagParserGenerator{
			testGenerator: testGenerator{tplugin: tplugin("mw")},
			err:           fmt.Errorf("error"),
		}

		pls := plugin.Plugins{
			erp, // A plugin that errors
			sec,
		}

		g.Receive(pls)
		g.ParseFlags([]string{"tg", "foo", "bar"})

		if len(sec.args) != 2 {
			t.Fatalf("expected generator to receive the flags")
		}

		if len(erp.args) != 2 {
			t.Fatalf("expected generator to receive the flags")
		}
	})
}

func TestGenerateMain(t *testing.T) {
	var g = generate.Command

	t.Run("call generate with invalid generator arg", func(t *testing.T) {
		out := &bytes.Buffer{}
		g.IO = &clio.IO{
			Out: out,
			Err: out,
		}

		g.Receive(plugin.Plugins{
			testGenerator{
				tplugin: tplugin("mw"),
			},
		})

		err := g.Main(context.Background(), "", []string{"xxxx"})
		if err != nil {
			t.Fatalf("should not error")
		}

		if !strings.Contains(out.String(), "No generator found for 'xxxx'") {
			t.Fatalf("should say it was not found")
		}
	})

	t.Run("call generator empty", func(t *testing.T) {
		out := &bytes.Buffer{}
		g.IO = &clio.IO{
			Out: out,
			Err: out,
		}

		err := g.Main(context.Background(), "", []string{"xxxx"})
		if err != nil {
			t.Fatalf("should not error")
		}

		if !strings.Contains(out.String(), "No generator found for 'xxxx'") {
			t.Fatalf("should say it was not found")
		}
	})

	t.Run("call existing generator", func(t *testing.T) {
		out := &bytes.Buffer{}
		g.IO = &clio.IO{
			Out: out,
			Err: out,
		}

		g.Receive(plugin.Plugins{
			testGenerator{
				tplugin: tplugin("mw"),
				IO: &clio.IO{
					Out: out,
					Err: out,
				},
			},
		})

		err := g.Main(context.Background(), "", []string{"mw"})
		if err != nil {
			t.Fatalf("should not error")
		}

		exp := `generating: mw`
		if !strings.Contains(out.String(), exp) {
			t.Fatalf("%s, should contain %s", out.String(), exp)
		}
	})
}
