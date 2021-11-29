package build

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"golang.org/x/mod/modfile"
)

func archivedAssets(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	mod, err := os.ReadFile(filepath.Join(opts.App.Root, "go.mod"))
	if err != nil {
		return g, err
	}

	mf, err := modfile.Parse("go.mod", mod, nil)
	if err != nil {
		return g, err
	}

	target := filepath.Join(filepath.Dir(opts.App.Bin), "assets.zip")
	source := filepath.Join(opts.App.Root, "public", "assets")

	g.RunFn(func(r *genny.Runner) error {
		bb := &bytes.Buffer{}
		archive := zip.NewWriter(bb)
		defer archive.Close()

		for _, f := range r.Disk.Files() {
			rel, err := filepath.Rel(source, f.Name())
			if err != nil {
				return err
			}

			if strings.HasPrefix(rel, "..") {
				continue
			}

			file, ok := f.(fs.File)
			if !ok {
				return fmt.Errorf("cannot process file %s", f.Name())
			}

			info, err := file.Stat()
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name = rel
			header.Method = zip.Deflate

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		// We need to close the archive before passing the buffer to genny, otherwise the zip
		// will be corrupted.
		archive.Close()
		if err := r.File(genny.NewFile(target, bb)); err != nil {
			return err
		}
		opts.keep.Store(target, struct{}{})
		return nil
	})

	g.RunFn(func(r *genny.Runner) error {
		f, err := r.FindFile("actions/app.go")
		if err != nil {
			return err
		}

		// public folder import is based on the module name
		// we will need to comment this when extracting assets
		// as it will not be used.
		pim := fmt.Sprintf(`"%v/public"`, mf.Module.Mod)

		opts.rollback.Store(f.Name(), f.String())
		body := strings.Replace(f.String(), `app.ServeFiles("/assets"`, `// app.ServeFiles("/assets"`, 1)
		body = strings.Replace(body, `app.ServeFiles("/"`, `// app.ServeFiles("/"`, 1)
		body = strings.Replace(body, pim, "//"+pim, 1)
		body = strings.Replace(body, `"net/http"`, `// "net/http"`, 1)

		return r.File(genny.NewFileS(f.Name(), body))
	})

	return g, nil
}
