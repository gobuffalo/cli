package pop

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	popfix "github.com/gobuffalo/pop/v6/fix"
)

var Fix = &fix{
	migrationsPath: "migrations",
}

type fix struct {
	migrationsPath string
}

func (c fix) Name() string {
	return "fix"
}

func (c fix) Usage() string {
	return "buffalo db fix"
}

func (c fix) HelpText() string {
	return "Brings pop, soda, and fizz files in line with the latest APIs"
}

func (c fix) PopMain(ctx context.Context, pwd string, args []string) error {
	return filepath.Walk(c.migrationsPath, func(path string, info os.FileInfo, _ error) error {
		if info == nil {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".fizz" {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		defer f.Close()
		bb := bytes.NewBuffer([]byte{})

		err = popfix.Fizz(f, bb)
		if err != nil {
			return err
		}

		// Fizz func does not write if there is not change.
		if bb.String() == "" {
			return nil
		}

		return ioutil.WriteFile(path, bb.Bytes(), 0644)
	})
}
