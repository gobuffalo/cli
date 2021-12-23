package fix

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
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
			b, err = addImport(path, b, fmt.Sprintf("%s/public", opts.App.PackagePkg))
			if err != nil {
				return err
			}

			rx := regexp.MustCompile("AssetsBox:.*,")
			b = rx.ReplaceAll(b, []byte("AssetsFS: public.FS(),"))

			rx = regexp.MustCompile(`^.*assetsBox.*=.*packr\.New.*$`)
			b = rx.ReplaceAll(b, []byte(""))
		}

		if bytes.Contains(b, []byte("TemplatesBox")) {
			b, err = addImport(path, b, fmt.Sprintf("%s/templates", opts.App.PackagePkg))
			if err != nil {
				return err
			}

			rx := regexp.MustCompile("TemplatesBox:.*,")
			b = rx.ReplaceAll(b, []byte("TemplatesFS: templates.FS(),"))
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

func addImport(path string, src []byte, importSpec string) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, src, 0)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(f.Decls); i++ {
		d := f.Decls[i]

		switch dd := d.(type) {
		case *ast.GenDecl:
			if dd.Tok == token.IMPORT {
				// Add the new import
				iSpec := &ast.ImportSpec{Path: &ast.BasicLit{Value: strconv.Quote(importSpec)}}
				dd.Specs = append(dd.Specs, iSpec)
			}
		default:
			// no action
		}
	}

	bb := &bytes.Buffer{}
	err = printer.Fprint(bb, fset, f)
	return bb.Bytes(), err
}
