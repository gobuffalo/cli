package fix

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
)

// UpdatePlushTemplates will update foo.html templates to foo.plush.html templates
func UpdatePlushTemplates(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		fmt.Println("~~~ Adding .plush extension to .html/.js/.md files ~~~")
		err := walkDisk(r.Disk, "templates", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			dir := filepath.Dir(path)
			base := filepath.Base(path)

			var exts []string
			ext := filepath.Ext(base)
			switch ext {
			case ".html":
			case ".js":
			case ".md":
			default:
				return nil
			}

			for len(ext) != 0 {
				if ext == ".plush" || ext == ".fizz" {
					return nil
				}
				exts = append([]string{ext}, exts...)
				base = strings.TrimSuffix(base, ext)
				ext = filepath.Ext(base)
			}
			exts = append([]string{base, ".plush"}, exts...)
			pathNew := filepath.Join(dir, strings.Join(exts, ""))

			fo, err := r.FindFile(path)
			if err != nil {
				return err
			}

			fn := genny.NewFile(pathNew, fo)
			if err := r.File(fn); err != nil {
				return err
			}
			return r.Delete(path)
		})
		return err
	}
}
