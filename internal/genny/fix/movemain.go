package fix

import (
	"github.com/gobuffalo/genny/v2"
)

// MoveMain will move the main.go from the root folder into
// cmd/app/main.go
func MoveMain(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		// fmt.Println("~~~ Checking main.go ~~~")
		// f, err := r.FindFile("main.go")
		// if err != nil {
		// 	return nil
		// }

		return nil
	}
}
