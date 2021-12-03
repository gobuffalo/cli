package build

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen/goimports"
)

func archivedAssets(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
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
				continue
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

		opts.rollback.Store(f.Name(), f.String())
		body := strings.Replace(f.String(), `app.ServeFiles("/assets"`, `// app.ServeFiles("/assets"`, 1)
		body = strings.Replace(body, `app.ServeFiles("/"`, `// app.ServeFiles("/"`, 1)
		f = genny.NewFileS(f.Name(), body)

		// run goimports after to ensure we remove unneeded imports
		// from actions/app.go
		content := bytes.NewBufferString("")
		gi := goimports.NewFromFiles(goimports.File{
			Name: f.Name(),
			In:   f,
			Out:  content,
		})

		if err := gi.Run(); err != nil {
			return err
		}

		return r.File(genny.NewFile(f.Name(), content))
	})

	return g, nil
}
