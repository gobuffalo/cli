package cli_test

import (
	"bytes"
	"context"
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/cli/cmd/cli"
)

//go:embed testdata/main.go.tmpl
var customMain []byte

func TestMain(t *testing.T) {

	t.Run("should render the general help if no command specified", func(t *testing.T) {
		a := cli.NewApp()
		bb := bytes.NewBuffer([]byte{})

		a.IO.Out = bb
		err := a.Main(context.TODO(), "", []string{})
		if err != nil {
			t.Fatalf("got error: %s", err)
		}

		if !strings.Contains(bb.String(), "Usage:") {
			t.Errorf("expected to contain 'Usage:'")
		}
	})

	t.Run("should render the general help if no command found", func(t *testing.T) {
		a := cli.NewApp()
		bb := bytes.NewBuffer([]byte{})

		a.IO.Out = bb
		err := a.Main(context.TODO(), "", []string{"no-command-with-this-name"})
		if err != nil {
			t.Fatalf("got error: %s", err)
		}

		if !strings.Contains(bb.String(), "Usage:") {
			t.Errorf("expected to contain 'Usage:'")
		}
	})

}

func TestCustomCLI(t *testing.T) {
	setupGlobalOverride := func() func() {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("could not get user home dir: %s", err)
		}

		path := filepath.Join(home, ".buffalo", "cmd")
		err = os.MkdirAll(path, 0755)
		if err != nil {
			t.Fatalf("could not create folder: %s", err)
		}

		err = os.WriteFile(filepath.Join(path, "main.go"), customMain, 0644)
		if err != nil {
			t.Fatalf("could not create file: %s", err)
		}

		return func() {
			os.RemoveAll(path)
		}
	}

	setupProjectOverride := func(pwd string) func() {
		path := filepath.Join(pwd, "cmd", "buffalo")
		err := os.MkdirAll(path, 0755)
		if err != nil {
			t.Fatalf("could not create folder: %s", err)
		}

		err = os.WriteFile(filepath.Join(path, "main.go"), customMain, 0644)
		if err != nil {
			t.Fatalf("could not create file: %s", err)
		}

		return func() {
			os.RemoveAll(path)
		}
	}

	t.Run("No overriders found", func(t *testing.T) {
		pwd := t.TempDir()

		bb := bytes.NewBuffer([]byte{})
		app := cli.NewApp()
		app.Out = bb

		err := app.Main(context.TODO(), pwd, []string{"plugins"})
		if err != nil {
			t.Fatal(err)
		}

		if strings.Contains(bb.String(), "[Info] Running CLI in `cmd/buffalo`") {
			t.Errorf("expected to not contain '[Info] Running CLI in `cmd/buffalo`'")
		}
	})

	t.Run("With project overrider", func(t *testing.T) {
		pwd := t.TempDir()
		cleanup := setupProjectOverride(pwd)
		t.Cleanup(cleanup)

		bb := bytes.NewBuffer([]byte{})
		app := cli.NewWithOverriders()

		app.IO.Out = bb
		err := app.Main(context.TODO(), pwd, []string{"plugins"})
		if err != nil {
			t.Fatalf("got error: %s", err)
		}

		if !strings.Contains(bb.String(), "[Info] Running CLI in `cmd/buffalo`") {
			t.Errorf("expected to contain '[Info] Running CLI in `cmd/buffalo`'")
		}
	})

	t.Run("With home override", func(t *testing.T) {
		t.Cleanup(setupGlobalOverride())

		bb := bytes.NewBuffer([]byte{})
		app := cli.NewWithOverriders()

		app.IO.Out = bb
		err := app.Main(context.TODO(), "", []string{"plugins"})
		if err != nil {
			t.Fatalf("got error: %s", err)
		}

		home, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("could not get user home dir: %s", err)
		}

		if !strings.Contains(bb.String(), "[Info] Running CLI in `"+home+"/buffalo`") {
			t.Errorf("expected to contain '[Info] Running CLI in `" + home + "/buffalo`'")
		}
	})

	t.Run("With both overriders should pick project", func(t *testing.T) {
		t.Cleanup(setupGlobalOverride())

		pwd := t.TempDir()
		t.Cleanup(setupProjectOverride(pwd))

		bb := bytes.NewBuffer([]byte{})
		app := cli.NewWithOverriders()

		app.IO.Out = bb
		err := app.Main(context.TODO(), pwd, []string{"plugins"})
		if err != nil {
			t.Fatalf("got error: %s", err)
		}

		if !strings.Contains(bb.String(), "[Info] Running CLI in `cmd/buffalo`") {
			t.Errorf("expected to contain '[Info] Running CLI in `cmd/buffalo`'")
		}
	})
}
