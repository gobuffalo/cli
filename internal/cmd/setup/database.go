package setup

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/gobuffalo/meta"
)

func databaseCheck(app meta.App) error {
	if !app.WithPop {
		return nil
	}
	for _, check := range []setupCheck{dbCreateCheck, dbMigrateCheck, dbSeedCheck} {
		err := check(app)
		if err != nil {
			return err
		}
	}
	return nil
}

func dbCreateCheck(meta.App) error {
	if setupOptions.dropDatabases {
		err := run(exec.Command("buffalo", "pop", "drop", "-a"))
		if err != nil {
			return fmt.Errorf("We encountered an error when trying to drop your application's databases. Please check to make sure that your database server is running and that the username and passwords found in the database.yml are properly configured and set up on your database server.\n %s", err)
		}
	}
	err := run(exec.Command("buffalo", "pop", "create", "-a"))
	if err != nil {
		return fmt.Errorf("We encountered an error when trying to create your application's databases. Please check to make sure that your database server is running and that the username and passwords found in the database.yml are properly configured and set up on your database server.\n %s", err)
	}
	return nil
}

func dbMigrateCheck(meta.App) error {
	err := run(exec.Command("buffalo", "pop", "migrate"))
	if err != nil {
		return fmt.Errorf("We encountered the following error when trying to migrate your database:\n%s", err)
	}
	return nil
}

func dbSeedCheck(meta.App) error {
	cmd := exec.Command("buffalo", "t", "list")
	out, err := cmd.Output()
	if err != nil {
		// no tasks configured, so return
		return nil
	}
	if bytes.Contains(out, []byte("db:seed")) {
		err := run(exec.Command("buffalo", "task", "db:seed"))
		if err != nil {
			return fmt.Errorf("We encountered the following error when trying to seed your database:\n%s", err)
		}
	}
	return nil
}
