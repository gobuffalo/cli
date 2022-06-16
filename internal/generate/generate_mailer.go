package generate

import (
	"context"
	"flag"

	"github.com/gobuffalo/cli/internal/genny/mail"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/gobuffalo/meta"
)

// func init() {
// 	MailCmd.Flags().BoolVarP(&mailOptions.dryRun, "dry-run", "d", false, "dry run of the generator")
// 	MailCmd.Flags().BoolVar(&mailOptions.SkipInit, "skip-init", false, "skip initializing mailers/")
// }

var MailerGenerator = &mailerGenerator{
	flagSet: flag.NewFlagSet("action", flag.ContinueOnError),
	options: &mail.Options{},
}

type mailerGenerator struct {
	app     meta.App
	options *mail.Options
	flagSet *flag.FlagSet

	dryRun   bool
	SkipInit bool
}

func (ag mailerGenerator) Name() string {
	return "mailer"
}

func (ag mailerGenerator) HelpText() string {
	return "Generates a new mailer for Buffalo"
}

func (ag mailerGenerator) Aliases() []string {
	return []string{"m"}
}

func (ag mailerGenerator) Generate(ctx context.Context, pwd string, args []string) error {
	ag.options.App = meta.New(".")

	ag.options.Name = name.New(args[0])
	gg, err := mail.New(ag.options)
	if err != nil {
		return err
	}

	run := genny.WetRunner(context.Background())
	if ag.dryRun {
		run = genny.DryRunner(context.Background())
	}

	g, err := gogen.Fmt(ag.app.Root)
	if err != nil {
		return err
	}

	if err := run.With(g); err != nil {
		return err
	}

	gg.With(run)
	return run.Run()
}
