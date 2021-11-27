package build

import (
	"bytes"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/cli/internal/genny/assets/webpack"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny/v2"
)

func assets(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	if opts.App.WithNodeJs || opts.App.WithWebpack {
		if opts.CleanAssets {
			g.RunFn(func(r *genny.Runner) error {
				return r.Delete(filepath.Join(opts.App.Root, "public", "assets"))
			})
		}
		g.RunFn(func(r *genny.Runner) error {
			r.Logger.Debugf("setting NODE_ENV = %s", opts.Environment)
			return envy.MustSet("NODE_ENV", opts.Environment)
		})
		g.RunFn(func(r *genny.Runner) error {
			tool := "yarnpkg"
			if !opts.App.WithYarn {
				tool = "npm"
			}

			c := exec.CommandContext(r.Context, tool, "run", "build")
			if _, err := opts.App.NodeScript("build"); err != nil {
				// Fallback on legacy runner
				c = exec.CommandContext(r.Context, webpack.BinPath)
			}

			bb := &bytes.Buffer{}
			c.Stdout = bb
			c.Stderr = bb

			if err := r.Exec(c); err != nil {
				r.Logger.Error(bb.String())
				return err
			}
			return nil
		})
	}

	if opts.ExtractAssets && opts.WithAssets {
		// mount the archived assets generator
		aa, err := archivedAssets(opts)
		if err != nil {
			return g, err
		}
		g.Merge(aa)
	}

	return g, nil
}
