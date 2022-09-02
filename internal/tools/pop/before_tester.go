package pop

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"

	"github.com/gobuffalo/pop/v6"
)

// BeforeTester will load the schema
// from the migrations folder if found. This is useful to speed
// tests from running migrations on each iteration.
type BeforeTester struct {
	forceMigrations bool
}

func (ls BeforeTester) Name() string {
	return "pop/setup-database"
}

func (ls BeforeTester) HelpText() string {
	return "Sets the database up before the tests run."
}

func (ls *BeforeTester) ParseFlags(args []string) (*flag.FlagSet, error) {
	fls := flag.NewFlagSet("pop/setup-database", flag.ContinueOnError)
	fls.BoolVar(&ls.forceMigrations, "force-migrations", false, "force migrations to run")
	_ = fls.Parse(args)

	return fls, nil
}

func (ls *BeforeTester) BeforeTest(ctx context.Context, pwd string, args []string) error {
	test, err := pop.Connect("test")

	// If its that it could not find the database config
	// file we should simply return  as the database cannot
	// be reset.
	if errors.Is(err, pop.ErrConfigFileNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	// drop the test db:
	if err := test.Dialect.DropDB(); err != nil {
		// not an error, since the database will be created in the next step anyway
		fmt.Println("[Info] no test database to drop")
	}

	// create the test db:
	err = test.Dialect.CreateDB()
	if err != nil {
		return err
	}

	if ls.forceMigrations {
		fm, err := pop.NewFileMigrator("./migrations", test)
		if err != nil {
			return err
		}

		if err := fm.Up(); err != nil {
			return err
		}

		return nil
	}

	schema := findSchema()
	if schema == nil {
		return nil
	}

	err = test.Dialect.LoadSchema(schema)
	if err != nil {
		return err
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
