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

func TestCustomCLI(t *testing.T) {

	t.Run("Just the app", func(t *testing.T) {
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

	t.Run("With custom CLI", func(t *testing.T) {
		pwd := t.TempDir()
		err := os.MkdirAll(filepath.Join(pwd, "cmd", "buffalo"), 0755)
		if err != nil {
			t.Fatalf("could not create folder: %s", err)
		}

		err = os.WriteFile(filepath.Join(pwd, "cmd", "buffalo", "main.go"), customMain, 0644)
		if err != nil {
			t.Fatalf("could not create file: %s", err)
		}

		bb := bytes.NewBuffer([]byte{})
		app := cli.DefaultApp

		app.IO.Out = bb
		err = app.Main(context.TODO(), pwd, []string{"plugins"})
		if err != nil {
			t.Fatalf("got error: %s", err)
		}

		t.Log(bb.String())

		if !strings.Contains(bb.String(), "[Info] Running CLI in `cmd/buffalo`") {
			t.Errorf("expected to contain '[Info] Running CLI in `cmd/buffalo`'")
		}
	})

}
