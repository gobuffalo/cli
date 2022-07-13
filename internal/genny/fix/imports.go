package fix

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gobuffalo/genny/v2"
	"golang.org/x/tools/go/ast/astutil"
)

var replace = map[string]string{
	"github.com/gobuffalo/buffalo-plugins":          "github.com/gobuffalo/cli/internal/plugins",
	"github.com/gobuffalo/buffalo-pop/":             "github.com/gobuffalo/buffalo-pop/v3",
	"github.com/gobuffalo/buffalo-pop/v2/":          "github.com/gobuffalo/buffalo-pop/v3",
	"github.com/gobuffalo/buffalo-pop/pop/popmw":    "github.com/gobuffalo/buffalo-pop/v3/pop/popmw",
	"github.com/gobuffalo/buffalo-pop/v2/pop/popmw": "github.com/gobuffalo/buffalo-pop/v3/pop/popmw",
	"github.com/gobuffalo/genny":                    "github.com/gobuffalo/genny/v2",
	"github.com/gobuffalo/mw-i18n":                  "github.com/gobuffalo/mw-i18n/v2",
	"github.com/gobuffalo/packr/v2":                 "",
	"github.com/gobuffalo/plush":                    "github.com/gobuffalo/plush/v4",
	"github.com/gobuffalo/pop":                      "github.com/gobuffalo/pop/v6",
	"github.com/gobuffalo/pop/v5":                   "github.com/gobuffalo/pop/v6",
	"github.com/gobuffalo/pop/nulls":                "github.com/gobuffalo/nulls",
	"github.com/gobuffalo/uuid":                     "github.com/gofrs/uuid",
	"github.com/gobuffalo/validate":                 "github.com/gobuffalo/validate/v3",
	"github.com/gobuffalo/validate/validators":      "github.com/gobuffalo/validate/v3/validators",
	"github.com/gobuffalo/suite":                    "github.com/gobuffalo/suite/v4",
	"github.com/markbates/grift":                    "github.com/gobuffalo/grift",
	"github.com/markbates/grift/grift":              "github.com/gobuffalo/grift/grift",
	"github.com/markbates/pop":                      "github.com/gobuffalo/pop/v6",
	"github.com/markbates/validate":                 "github.com/gobuffalo/validate/v3",
	"github.com/markbates/willie":                   "github.com/gobuffalo/httptest",
	"github.com/satori/go.uuid":                     "github.com/gofrs/uuid",
	"github.com/shurcooL/github_flavored_markdown":  "github.com/gobuffalo/github_flavored_markdown",
}

// ReplaceOldImports walks all the .go files in an application
// It will then attempt to convert any old import paths to any new import paths
// used by this version Buffalo.
func ReplaceOldImports(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		fmt.Println("~~~ Rewriting Imports ~~~")
		return walkDisk(r.Disk, ".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".go" {
				return nil
			}

			f, err := r.FindFile(path)
			if err != nil {
				return err
			}
			if err := rewriteFile(f); err != nil {
				return err
			}
			return r.File(f)
		})
	}
}

func rewriteFile(file genny.File) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file.Name(), file.String(), parser.ParseComments)
	if err != nil {
		return err
	}

	names := make(map[string]string)
	for _, imports := range astutil.Imports(fset, f) {
		for _, imp := range imports {
			i, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				return err
			}
			if imp.Name != nil {
				names[i] = imp.Name.Name
			} else {
				names[i] = ""
			}
		}
	}

	changed := false
	for key, value := range replace {
		name, ok := names[key]
		if !ok {
			continue
		}

		astutil.DeleteNamedImport(fset, f, name, key)
		if value != "" {
			astutil.AddImport(fset, f, value)
		}

		changed = true
	}

	// if no change occurred, then we don't need to write to disk, just return.
	if !changed {
		return nil
	}

	// since the imports changed, resort them.
	ast.SortImports(fset, f)

	bb := &bytes.Buffer{}
	if err := printer.Fprint(bb, fset, f); err != nil {
		return err
	}

	_, err = file.Write(bb.Bytes())
	return err
}
