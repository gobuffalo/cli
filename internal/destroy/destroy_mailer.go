package destroy

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/flect"
)

var MailerDestroyer = &mailerDestroyer{}

type mailerDestroyer struct {
	yesToAll bool
}

func (ad mailerDestroyer) Name() string {
	return "mailer"
}

func (ad mailerDestroyer) Usage() string {
	return "destroy [flags] mailer [name]"
}

func (ad mailerDestroyer) HelpText() string {
	return "Destroy mailer files"
}

func (ad mailerDestroyer) Aliases() []string {
	return []string{"l"}
}

func (ad *mailerDestroyer) PreConfirm() {
	ad.yesToAll = true
}

func (ad *mailerDestroyer) Destroy(ctx context.Context, pwd string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you need to provide a valid mailer name in order to destroy it")
	}

	if !ad.yesToAll && !confirm("Want to remove mailer? (y/N)") {
		return nil
	}

	name := args[0]
	mailerFileName := flect.Singularize(flect.Underscore(name))

	files := []string{
		filepath.Join("mailers", fmt.Sprintf("%v.go", mailerFileName)),
		filepath.Join("templates/mail", fmt.Sprintf("%v.html", mailerFileName)),
		filepath.Join("templates/mail", fmt.Sprintf("%v.plush.html", mailerFileName)),
	}

	for _, f := range files {
		// TODO: Handle error here.
		os.Remove(f)
		fmt.Printf("- Deleted %v\n", f)
	}

	return nil
}
