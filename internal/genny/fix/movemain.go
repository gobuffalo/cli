package fix

import (
	"fmt"
	"io"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/refresh/refresh"
	"gopkg.in/yaml.v2"
)

// MoveMain will move the main.go from the root folder into
// cmd/app/main.go
func MoveMain(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		// If there is a main in the cmd/app folder, we don't need to do anything
		_, err := r.FindFile("cmd/app/main.go")
		if err == nil {
			return nil
		}

		fmt.Println("~~~ Moving main.go ~~~")
		// If there is a main in the root folder, we need to move it
		f, err := r.FindFile("main.go")
		if err != nil {
			// There is no file to move to move\
			r.Logger.Info("No main.go found")
			return nil
		}

		nf := genny.NewFileS("cmd/app/main.go", f.String())
		err = r.Disk.Delete("main.go")
		if err != nil {
			return fmt.Errorf("could not delete main.go: %w", err)
		}

		return r.File(nf)
	}
}

// Refresh will update the buffalo.dev.yml to build
// the main at ./cmd/app instead of the root folder.
func Refresh(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		fmt.Println("~~~ Checking for .buffalo.dev.yml ~~~")
		f, err := r.FindFile(".buffalo.dev.yml")
		if err != nil {
			return err
		}

		data, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		c := refresh.Configuration{}
		err = yaml.Unmarshal(data, &c)
		if err != nil {
			return err
		}

		if c.BuildTargetPath != "" {
			return nil
		}

		c.BuildTargetPath = "./cmd/app"
		data, err = yaml.Marshal(&c)
		if err != nil {
			return err
		}

		_, err = f.Write(data)
		if err != nil {
			return err
		}

		return r.File(f)
	}
}
