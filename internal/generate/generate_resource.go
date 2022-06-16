package generate

import (
	"context"
	"flag"
	"fmt"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/cli/internal/genny/resource"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/meta"
)

// func init() {
// 	ResourceCmd.Flags().BoolVarP(&resourceOptions.SkipMigration, "skip-migration", "s", false, "tells resource generator not-to add model migration")
// 	ResourceCmd.Flags().BoolVarP(&resourceOptions.SkipModel, "skip-model", "", false, "tells resource generator not to generate model nor migrations")
// 	ResourceCmd.Flags().BoolVarP(&resourceOptions.SkipTemplates, "skip-templates", "", false, "tells resource generator not to generate templates for the resource")
// 	ResourceCmd.Flags().StringVarP(&resourceOptions.Model, "use-model", "", "", "tells resource generator to reference an existing model in generated code")
// 	ResourceCmd.Flags().StringVarP(&resourceOptions.Name, "name", "n", "", "allows to define a different model name for the resource being generated.")
// 	ResourceCmd.Flags().BoolVarP(&resourceOptions.DryRun, "dry-run", "d", false, "dry run")
// 	ResourceCmd.Flags().BoolVarP(&resourceOptions.Verbose, "verbose", "v", false, "verbosely print out the go get commands")
// }

var ResourceGenerator = &resourceGenerator{
	flagSet: flag.NewFlagSet("resource", flag.ContinueOnError),
	options: &resource.Options{},
}

type resourceGenerator struct {
	app     meta.App
	options *resource.Options
	flagSet *flag.FlagSet

	dryRun  bool
	verbose bool
}

func (ag resourceGenerator) Name() string {
	return "resource"
}

func (ag resourceGenerator) HelpText() string {
	return "Generate a new actions/resource file"
}

func (ag resourceGenerator) LongHelpText() string {
	return `Example Usage:
	
	$ buffalo g resource users
	Generates:
	
	- actions/users.go
	- actions/users_test.go
	- models/user.go
	- models/user_test.go
	- migrations/2016020216301234_create_users.up.fizz
	- migrations/2016020216301234_create_users.down.fizz
	
	$ buffalo g resource users --skip-migration
	Generates:
	
	- actions/users.go
	- actions/users_test.go
	- models/user.go
	- models/user_test.go
	
	$ buffalo g resource users --skip-model
	Generates:
	
	- actions/users.go
	- actions/users_test.go
	
	$ buffalo g resource users --use-model users
	Generates:
	
	- actions/users.go
	- actions/users_test.go`
}

func (ag resourceGenerator) Aliases() []string {
	return []string{"r"}
}

func (ag resourceGenerator) Generate(ctx context.Context, pwd string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must supply a name")
	}

	run := genny.WetRunner(ctx)
	if ag.dryRun {
		run = genny.DryRunner(ctx)
	}

	if ag.verbose {
		lg := logger.New(logger.DebugLevel)
		run.Logger = lg
	}

	if len(ag.options.Name) == 0 {
		ag.options.Name = args[0]
	}

	if len(args) > 1 {
		ats, err := attrs.ParseArgs(args[1:]...)
		if err != nil {
			return err
		}

		ag.options.Attrs = ats
	}

	if err := run.WithNew(resource.New(ag.options)); err != nil {
		return err
	}

	return run.Run()
}
