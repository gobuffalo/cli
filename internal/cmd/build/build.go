package build

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gobuffalo/cli/internal/genny/build"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/meta"
	"github.com/spf13/cobra"
)

var buildOptions = struct {
	*build.Options
	SkipAssets             bool
	SkipBuildDeps          bool
	Debug                  bool
	Tags                   string
	SkipTemplateValidation bool
	DryRun                 bool
	Verbose                bool
	bin                    string
	WithCustomModule       bool
}{
	Options: &build.Options{
		BuildTime: time.Now(),
	},
}

func runE(cmd *cobra.Command, args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	buildOptions.App = meta.New(pwd)
	if len(buildOptions.bin) > 0 {
		buildOptions.App.Bin = buildOptions.bin
	}
	if buildOptions.WithCustomModule {
		buildOptions.App.PackagePkg = buildOptions.App.Name.Original
		actionsPkg := strings.Split(buildOptions.App.ActionsPkg, "/")
		if len(actionsPkg) > 1 {
			buildOptions.App.ActionsPkg = buildOptions.App.Name.Original + "/" + actionsPkg[1]
		} else {
			buildOptions.App.ActionsPkg = buildOptions.App.Name.Original + "/" + actionsPkg[0]
		}
		modelsPkg := strings.Split(buildOptions.App.ModelsPkg, "/")
		if len(modelsPkg) > 1 {
			buildOptions.App.ModelsPkg = buildOptions.App.Name.Original + "/" + modelsPkg[1]
		} else {
			buildOptions.App.ModelsPkg = buildOptions.App.Name.Original + "/" + modelsPkg[0]
		}
		griftsPkg := strings.Split(buildOptions.App.GriftsPkg, "/")
		if len(griftsPkg) > 1 {
			buildOptions.App.GriftsPkg = buildOptions.App.Name.Original + "/" + griftsPkg[1]
		} else {
			buildOptions.App.GriftsPkg = buildOptions.App.Name.Original + "/" + griftsPkg[0]
		}
	}

	buildOptions.Options.WithAssets = !buildOptions.SkipAssets
	buildOptions.Options.WithBuildDeps = !buildOptions.SkipBuildDeps

	run := genny.WetRunner(ctx)
	if buildOptions.DryRun {
		run = genny.DryRunner(ctx)
	}

	if buildOptions.Verbose || buildOptions.Debug {
		lg := logger.New(logger.DebugLevel)
		run.Logger = lg
		buildOptions.BuildFlags = append(buildOptions.BuildFlags, "-v")
	}

	opts := buildOptions.Options
	opts.BuildVersion = buildVersion(opts.BuildTime.Format(time.RFC3339))

	if buildOptions.Tags != "" {
		opts.Tags = append(opts.Tags, buildOptions.Tags)
	}

	if !buildOptions.SkipTemplateValidation {
		opts.TemplateValidators = append(opts.TemplateValidators, build.PlushValidator, build.GoTemplateValidator)
	}

	if cmd.CalledAs() == "install" {
		opts.GoCommand = "install"
	}
	clean := build.Cleanup(opts)
	// defer clean(run)
	defer func() {
		if err := clean(run); err != nil {
			log.Fatalf("build:clean %s", err)
		}
	}()
	if err := run.WithNew(build.New(opts)); err != nil {
		return err
	}
	return run.Run()
}

func buildVersion(version string) string {
	vcs := buildOptions.VCS

	if len(vcs) == 0 {
		return version
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key != "vcs.revision" {
				continue
			}

			return setting.Value
		}
	}

	return version
}
