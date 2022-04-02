package fix

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/cli/internal/runtime"
	"github.com/gobuffalo/events"
	"github.com/gobuffalo/genny/v2"
)

func ask(q string) bool {
	fmt.Printf("? %s [y/n]\n", q)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	text = strings.ToLower(strings.TrimSpace(text))
	return text == "y" || text == "yes"
}

func printWarnings(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		if len(opts.warnings) == 0 {
			return nil
		}

		fmt.Println("\n\n----------------------------")
		fmt.Printf("!!! (%d) Warnings Were Found !!!\n\n", len(opts.warnings))
		for _, w := range opts.warnings {
			fmt.Printf("[WARNING]: %s\n", w)
		}
		return nil
	}
}

func tidyCmd() *exec.Cmd {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stderr = os.Stderr
	return cmd
}

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	fmt.Printf("! This updater will attempt to update your application to Buffalo version: %s\n", runtime.Version)
	if !opts.YesToAll && !ask("Do you wish to continue?") {
		fmt.Println("~~~ cancelling update ~~~")
		return g, nil
	}

	g.ErrorFn = func(err error) {
		events.EmitError(EvtFixStopErr, err, events.Payload{"opts": opts})
	}

	g.RunFn(func(r *genny.Runner) error {
		events.EmitPayload(EvtFixStart, events.Payload{"opts": opts})
		return nil
	})

	// ExceptFS will not scan any files that the path starts with given strings
	// excluded files include such as `.git/`, `.github/`, `.gitignore`, `.yarn/`, `.yarnrc`, etc.
	if err := g.ExceptFS(os.DirFS(opts.App.Root), []string{"node_modules", ".git", ".yarn"}, nil); err != nil {
		return nil, err
	}

	// replace old imports with new ones
	g.RunFn(ReplaceOldImports(opts))
	g.Command(tidyCmd())

	// check webpack.config.json and package.json for updates
	g.RunFn(WebpackCheck(opts))
	g.RunFn(PackageJSONCheck(opts))
	g.RunFn(AddPackageJSONScripts(opts))

	// install required tools
	g.RunFn(InstallTools(opts))
	g.Command(tidyCmd())

	// check for deprecations
	g.RunFn(DeprecationsCheck(opts))

	// add embed.go files
	g.RunFn(FixEmbed(opts))

	// fix Docker file
	g.RunFn(FixDocker(opts))

	g.RunFn(EncodeAppToml(opts))

	// update plugins
	g.RunFn(RemoveOldPlugins(opts))
	g.RunFn(CleanPluginCache)
	g.RunFn(ReinstallPlugins(opts))

	// update plush templates
	g.RunFn(UpdatePlushTemplates(opts))

	g.RunFn(MoveMain(opts))
	g.RunFn(Refresh(opts))

	// print all warnings that were captured
	g.RunFn(printWarnings(opts))

	g.RunFn(func(r *genny.Runner) error {
		events.EmitPayload(EvtFixStop, events.Payload{"opts": opts})
		return nil
	})
	return g, nil
}
