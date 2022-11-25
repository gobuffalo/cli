package test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/pop/v6"
)

func setupDatabase(args []string) error {
	if _, err := os.Stat("database.yml"); err != nil {
		return err
	}

	test, err := pop.Connect("test")
	if err != nil {
		return err
	}

	// drop the test db:
	if err := test.Dialect.DropDB(); err != nil {
		// not an error, since the database will be created in the next step anyway
		fmt.Println("INFO: no test database to drop.")
	}

	// create the test db:
	err = test.Dialect.CreateDB()
	if err != nil {
		return err
	}

	forceMigrations := strings.Contains(strings.Join(args, " "), "--force-migrations")
	if forceMigrations {
		fm, err := pop.NewFileMigrator("./migrations", test)
		if err != nil {
			return err
		}

		return fm.Up()
	}

	if schema := findSchema(); schema != nil {
		return test.Dialect.LoadSchema(schema)
	}

	return nil
}

func findSchema() io.Reader {
	if f, err := os.Open(filepath.Join("migrations", "schema.sql")); err == nil {
		return f
	}

	if dev, err := pop.Connect("development"); err == nil {
		schema := &bytes.Buffer{}
		if err = dev.Dialect.DumpSchema(schema); err == nil {
			return schema
		}
	}

	if test, err := pop.Connect("test"); err == nil {
		fm, err := pop.NewFileMigrator("./migrations", test)
		if err != nil {
			return nil
		}

		if err := fm.Up(); err == nil {
			return nil
		}
	}

	return nil
}
