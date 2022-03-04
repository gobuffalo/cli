package webpack

import (
	"embed"
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

//go:embed templates/* templates/assets/css/_buffalo.scss.tmpl
var templates embed.FS

// Templates used for generating webpack
// (exported mostly for the "fix" command)
func Templates() (fs.FS, error) {
	return fs.Sub(templates, "templates")
}

// BinPath is the path to the local install of webpack
var BinPath = func() string {
	s := filepath.Join("node_modules", ".bin", "webpack")
	if runtime.GOOS == "windows" {
		s += ".cmd"
	}
	return s
}()

// New generator for creating webpack asset files
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	g.RunFn(func(r *genny.Runner) error {
		if opts.App.WithYarn {
			if _, err := r.LookPath("yarnpkg"); err == nil {
				return nil
			}
			// If yarn is not installed, it still can be installed with npm.
		}
		if _, err := r.LookPath("npm"); err != nil {
			return fmt.Errorf("could not find npm executable")
		}
		return nil
	})

	temp, err := Templates()
	if err != nil {
		return g, err
	}

	if err := g.FS(temp); err != nil {
		return g, err
	}

	data := map[string]interface{}{
		"opts": opts,
	}
	t := gogen.TemplateTransformer(data, gogen.TemplateHelpers)
	g.Transformer(t)
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

	g.RunFn(func(r *genny.Runner) error {
		return installPkgs(r, opts)
	})

	return g, nil
}

func installPkgs(r *genny.Runner, opts *Options) error {
	command := "yarnpkg"
	args := []string{"install"}

	if !opts.App.WithYarn {
		command = "npm"
		args = []string{"install", "--no-progress", "--save"}
	} else {
		if err := installYarn(r); err != nil {
			return err
		}
		if err := r.Exec(exec.Command(command, []string{"set", "version", "berry"}...)); err != nil {
			return err
		}
	}

	c := exec.Command(command, args...)
	c.Stdout = yarnWriter{
		fn: r.Logger.Debug,
	}
	c.Stderr = yarnWriter{
		fn: r.Logger.Debug,
	}
	return r.Exec(c)
}

type yarnWriter struct {
	fn func(...interface{})
}

func (y yarnWriter) Write(p []byte) (int, error) {
	y.fn(string(p))
	return len(p), nil
}

func installYarn(r *genny.Runner) error {
	// if there's no yarn, install it!
	if _, err := r.LookPath("yarnpkg"); err == nil {
		return nil
	}
	yargs := []string{"install", "-g", "yarn"}
	return r.Exec(exec.Command("npm", yargs...))
}
