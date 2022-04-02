package fix

import (
	"fmt"
	"io/ioutil"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/refresh/refresh"
	"gopkg.in/yaml.v2"
)

// Refresh will update the buffalo.dev.yml to build
// the main at ./cmd/app instead of the root folder.
func Refresh(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		fmt.Println("~~~ Checking for .buffalo.dev.yml ~~~")
		f, err := r.FindFile(".buffalo.dev.yml")
		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(f)
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
