package destroy

import (
	"context"
	"fmt"

	"github.com/gobuffalo/flect"
)

var ActionDestroyer = &actionDestroyer{}

type actionDestroyer struct {
	yesToAll bool
}

func (ad actionDestroyer) Name() string {
	return "action"
}

func (ad actionDestroyer) HelpText() string {
	return "Destroys action files"
}

func (ad actionDestroyer) Usage() string {
	return "buffalo destroy [flags] action [name]"
}

func (ad actionDestroyer) Aliases() []string {
	return []string{"a"}
}

func (ad *actionDestroyer) PreConfirm() {
	ad.yesToAll = true
}

func (ad *actionDestroyer) Destroy(ctx context.Context, pwd string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you need to provide a valid action file name in order to destroy it")
	}

	name := args[0]

	// Generated actions keep the same name (not plural).
	fileName := flect.Underscore(name)
	return removeActions(ad.yesToAll, fileName)
}
