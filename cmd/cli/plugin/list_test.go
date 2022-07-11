package plugin_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

func TestList(t *testing.T) {
	cc := plugin.Plugins{
		testCommand("tt"),
	}

	out := bytes.NewBuffer([]byte{})
	plugin.List.Out = out

	plugin.List.Receive(cc)

	err := plugin.List.Main(context.Background(), "", []string{})
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range []string{
		"tt",
		"runs the tt thing basically",
		"Loaded default CLI plugins",
		"Plugins loaded (1)",
	} {
		if !strings.Contains(out.String(), v) {
			t.Fatalf("expected to not contain '%s'", v)
		}
	}

	if !strings.Contains(out.String(), "tt") {
		t.Fatal("expected to not contain 'Loaded default CLI plugins'")
	}
}
