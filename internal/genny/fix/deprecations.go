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
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"golang.org/x/tools/go/ast/astutil"
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

		err = walkDisk(r.Disk, "actions", actionsWalkFun(r.Disk, opts))
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

			rx = regexp.MustCompile(`^.*AssetsBox.*=.*$`)
			b = rx.ReplaceAll(b, []byte(""))
		}

		if bytes.Contains(b, []byte("TemplatesBox")) {
			b, err = addImport(path, b, fmt.Sprintf("%s/templates", opts.App.PackagePkg))
			if err != nil {
				return err
			}

			rx := regexp.MustCompile("TemplatesBox:.*,")
			b = rx.ReplaceAll(b, []byte("TemplatesFS: templates.FS(),"))

			rx = regexp.MustCompile(`^.*TemplatesBox.*=.*$`)
			b = rx.ReplaceAll(b, []byte(""))
		}

		b, err = format.Source(b)
		if err != nil {
			return err
		}

		_, err = f.Write(b)
		return err
	}
}

func addImport(path string, src []byte, importSpec string) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	astutil.AddImport(fset, f, importSpec)
	ast.SortImports(fset, f)

	bb := &bytes.Buffer{}
	err = printer.Fprint(bb, fset, f)
	return bb.Bytes(), err
}
