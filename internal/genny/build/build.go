package build

import (
	"embed"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/gobuffalo/cli/internal/runtime"
	"github.com/gobuffalo/events"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/plushgen"
	"github.com/gobuffalo/plush/v4"
)

//go:embed templates/*
var templates embed.FS

// New generator for building a Buffalo application
// This powers the `buffalo build` command and can be
// used to programatically build/compile Buffalo
// applications.
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}
	g.ErrorFn = func(err error) {
		events.EmitError(EvtBuildStopErr, err, events.Payload{"opts": opts})
	}

	g.RunFn(func(r *genny.Runner) error {
		events.EmitPayload(EvtBuildStart, events.Payload{"opts": opts})
		return nil
	})

	g.Transformer(genny.Dot())
	// TODO: workaround for 1.16, remove when we upgrade to 1.17 and rename "dot-*" files back to "-dot-*"
	g.Transformer(genny.NewTransformer("*", func(f genny.File) (genny.File, error) {
		name := f.Name()
		if strings.HasPrefix(name, "dot-") {
			name = strings.TrimPrefix(name, "dot-")
			name = "." + name
		}
		return genny.NewFile(name, f), nil
	}))
	g.Transformer(genny.Replace("/dot-", "/."))

	// validate templates
	templatesPath := opts.App.Root + "/templates"
	if _, err := os.Stat(templatesPath); err == nil {
		g.RunFn(ValidateTemplates(os.DirFS(templatesPath), opts.TemplateValidators))
	}

	// rename main() to originalMain()
	g.RunFn(transformMain(opts))

	// add any necessary templates for the build
	sub, err := fs.Sub(templates, "templates")
	if err != nil {
		return g, err
	}

	if err := g.FS(sub); err != nil {
		return g, err
	}

	// configure plush
	ctx := plush.NewContext()
	ctx.Set("opts", opts)
	ctx.Set("buildTime", opts.BuildTime.Format(time.RFC3339))
	ctx.Set("buildVersion", opts.BuildVersion)
	ctx.Set("buffaloVersion", runtime.Version)
	g.Transformer(plushgen.Transformer(ctx))

	// create the ./a pkg
	ag, err := apkg(opts)
	if err != nil {
		return g, err
	}
	g.Merge(ag)

	if opts.WithAssets {
		// mount the assets generator
		ag, err := assets(opts)
		if err != nil {
			return g, err
		}
		g.Merge(ag)
	}

	if opts.WithBuildDeps {
		// mount the build time dependency generator
		dg, err := buildDeps(opts)
		if err != nil {
			return g, err
		}
		g.Merge(dg)
	}

	// create the final go build command
	c, err := buildCmd(opts)
	if err != nil {
		return g, err
	}

	g.Command(c)
	g.RunFn(func(r *genny.Runner) error {
		events.EmitPayload(EvtBuildStop, events.Payload{"opts": opts})
		return nil
	})

	g.RunFn(Cleanup(opts))
	return g, nil
}
