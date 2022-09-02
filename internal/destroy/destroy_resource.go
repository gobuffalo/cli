package destroy

import (
	"context"
	"fmt"

	"github.com/gobuffalo/flect"
)

var ResourceDestroyer = &resourceDestroyer{}

type resourceDestroyer struct {
	yesToAll bool
}

func (ad resourceDestroyer) Name() string {
	return "resource"
}

func (ad resourceDestroyer) Usage() string {
	return "destroy [flags] resource [name]"
}

func (ad resourceDestroyer) HelpText() string {
	return "Destroy resource files"
}

func (ad resourceDestroyer) Aliases() []string {
	return []string{"r"}
}

func (ad *resourceDestroyer) PreConfirm() {
	ad.yesToAll = true
}

func (ad *resourceDestroyer) Destroy(ctx context.Context, pwd string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you need to provide a valid resource name in order to destroy it")
	}

	name := args[0]
	fileName := flect.Pluralize(flect.Underscore(name))

	removeTemplates(ad.yesToAll, fileName)
	if err := removeActions(ad.yesToAll, fileName); err != nil {
		return err
	}

	removeLocales(ad.yesToAll, fileName)
	removeModel(ad.yesToAll, fileName)
	removeMigrations(ad.yesToAll, fileName)

	return nil
}
