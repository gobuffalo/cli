package build

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
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

	ctx := context.Background()
	run := genny.WetRunner(ctx)
	if buildOptions.DryRun {
		run = genny.DryRunner(ctx)
	}

	_, err := exec.LookPath(vcs)
	if err != nil {
		run.Logger.Warnf("could not find %s; defaulting to version %s", vcs, version)
		return vcs
	}

	var cmd *exec.Cmd
	switch vcs {
	case "git":
		// If .git folder does not exist return default version
		if stat, err := os.Stat(".git"); err != nil || !stat.IsDir() {
			run.Logger.Warnf("could not find .git folder; defaulting to version %s", version)
			return version
		}

		cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
	case "bzr":
		cmd = exec.Command("bzr", "revno")
	default:
		run.Logger.Warnf("could not find %s; defaulting to version %s", vcs, version)
		return vcs
	}

	out := &bytes.Buffer{}
	cmd.Stdout = out
	run.WithRun(func(r *genny.Runner) error {
		return r.Exec(cmd)
	})

	if err := run.Run(); err != nil {
		run.Logger.Error(err)
		return version
	}

	if out.String() != "" {
		return strings.TrimSpace(out.String())
	}

	return version
}
