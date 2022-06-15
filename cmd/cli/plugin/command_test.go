package plugin_test

import (
	"testing"
)

func TestCommandFind(t *testing.T) {
	cc := Commands{
		testCommand("tt"),
	}

	c := cc.Find("tt")
	if c == nil {
		t.Fatalf("did not find test command by its name: tt")
	}

	c = cc.Find("tc")
	if c == nil {
		t.Fatalf("did not find test command by its alias: tc")
	}

}
