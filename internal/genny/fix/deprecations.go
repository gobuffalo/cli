package fix

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
)

// DeprecationsCheck will either log, or fix, deprecated items in the application
func DeprecationsCheck(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		fmt.Println("~~~ Checking for deprecations ~~~")
		f, err := r.FindFile("main.go")
		if err != nil {
			return err
		}

		if strings.Contains(f.String(), "app.Start") {
			opts.warnings = append(opts.warnings, "app.Start has been removed in v0.11.0. Use app.Serve Instead. [main.go]")
		}

		err = walkDisk(r.Disk, filepath.Join(opts.App.Root, "actions"), actionsWalkFun(r.Disk, opts))
		// TODO: add other folders to check
		return err
	}
}

func actionsWalkFun(r *genny.Disk, opts *Options) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		if err != nil {
			return err
		}

		f, err := r.Find(path)
		if err != nil {
			return err
		}
		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		if bytes.Contains(b, []byte("AssetsBox")) {
			b = bytes.Replace(b, []byte("AssetsBox:"), []byte("AssetsFS:"), -1)
		}

		if bytes.Contains(b, []byte("TemplatesBox")) {
			b = bytes.Replace(b, []byte("TemplatesBox:"), []byte("TemplatesFS:"), -1)
		}

		b, err = format.Source(b)
		if err != nil {
			return err
		}

		_, err = f.Write(b)
		return err
	}
}

func walkDisk(disk *genny.Disk, root string, walkFun filepath.WalkFunc) error {
	for _, f := range disk.Files() {
		rel, err := filepath.Rel(root, f.Name())
		if err != nil {
			err = walkFun(f.Name(), nil, fmt.Errorf("cannot process file %s", f.Name()))
			if err != nil {
				return err
			}
		}

		if strings.HasPrefix(rel, "..") {
			continue
		}

		file, ok := f.(fs.File)
		if !ok {
			return fmt.Errorf("cannot process file %s", f.Name())
		}

		info, err := file.Stat()

		err = walkFun(f.Name(), info, err)
		if err != nil {
			return err
		}
	}

	return nil
}
