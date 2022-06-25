package build

import (
	"bytes"
	"context"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gobuffalo/cli/cmd/cli/clio"
	"github.com/gobuffalo/cli/internal/genny/build"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/meta"
	"github.com/markbates/sigtx"
)

var Command = &command{
	options: &build.Options{
		BuildTime: time.Now(),
	},
}

type command struct {
	options *build.Options
	flagSet *flag.FlagSet

	skipAssets             bool
	skipBuildDeps          bool
	debug                  bool
	skipTemplateValidation bool
	dryRun                 bool
	verbose                bool
	bin                    string

	tags       clio.StringSliceVar
	buildFlags clio.StringSliceVar
}

func (c command) Name() string {
	return "build"
}

func (c command) HelpText() string {
	return "Build the application binary, including bundling of webpack assets"
}

func (c command) Aliases() []string {
	return []string{"b", "bill", "install"}
}

func (c *command) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(io.Discard)
	}

	c.flagSet.StringVar(&c.bin, "output", c.bin, "set the name of the binary")

	c.flagSet.BoolVar(&c.options.ExtractAssets, "extract-assets", false, "extract the assets and put them in a distinct archive")
	c.flagSet.BoolVar(&c.skipAssets, "skip-assets", false, "skip running webpack and building assets")
	c.flagSet.BoolVar(&c.skipBuildDeps, "skip-build-deps", false, "skip building dependencies")
	c.flagSet.BoolVar(&c.options.Static, "static", false, "build a static binary using  --ldflags '-linkmode external -extldflags \"-static\"'")
	c.flagSet.StringVar(&c.options.LDFlags, "ldflags", "", "set any ldflags to be passed to the go build")
	c.flagSet.BoolVar(&c.verbose, "verbose", false, "print debugging information")
	c.flagSet.BoolVar(&c.dryRun, "dry-run", false, "runs the build 'dry'")
	c.flagSet.BoolVar(&c.skipTemplateValidation, "skip-template-validation", false, "skip validating templates")
	c.flagSet.BoolVar(&c.options.CleanAssets, "clean-assets", false, "will delete public/assets before calling webpack")
	c.flagSet.StringVar(&c.options.Environment, "environment", "development", "set the environment for the binary")
	c.flagSet.StringVar(&c.options.Mod, "mod", "", "-mod flag for go build")

	c.flagSet.Var(&c.tags, "tags", "compile with specific build tags")
	c.flagSet.Var(&c.buildFlags, "build-flags", "Additional comma-separated build flags to feed to go build")

	return c.flagSet, nil
}

func (c *command) Main(ctx context.Context, pwd string, args []string) error {
	ctx, cancel := sigtx.WithCancel(ctx, os.Interrupt)
	defer cancel()

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	c.options.App = meta.New(pwd)
	if len(c.bin) > 0 {
		c.options.App.Bin = c.bin
	}

	c.options.WithAssets = !c.skipAssets
	c.options.WithBuildDeps = !c.skipBuildDeps

	run := genny.WetRunner(ctx)
	if c.dryRun {
		run = genny.DryRunner(ctx)
	}

	if c.verbose || c.debug {
		run.Logger = logger.New(logger.DebugLevel)
		c.buildFlags = append(c.buildFlags, "-v")
	}

	c.options.BuildVersion = c.buildVersion(c.options.BuildTime.Format(time.RFC3339))

	if len(c.tags) > 0 {
		c.options.Tags = meta.BuildTags(c.tags)
	}

	if !c.skipTemplateValidation {
		c.options.TemplateValidators = append(
			c.options.TemplateValidators,
			build.PlushValidator,
			build.GoTemplateValidator,
		)
	}

	// if cmd.CalledAs() == "install" {
	// 	opts.GoCommand = "install"
	// }

	clean := build.Cleanup(c.options)
	defer func() {
		if err := clean(run); err != nil {
			log.Fatalf("build:clean %s", err)
		}
	}()

	if err := run.WithNew(build.New(c.options)); err != nil {
		return err
	}

	return run.Run()
}

// TODO: This needs to be pluginized
func (c command) buildVersion(version string) string {
	vcs := c.options.VCS

	if len(vcs) == 0 {
		return version
	}

	ctx := context.Background()
	run := genny.WetRunner(ctx)
	if c.dryRun {
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
