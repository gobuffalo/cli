package fix

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
)

// DeprecationsCheck will either log, or fix, deprecated items in the application
func DeprecationsCheck(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		fmt.Println("~~~ Checking for deprecations ~~~")
		b, err := os.ReadFile("main.go")
		if err != nil {
			return err
		}

		if bytes.Contains(b, []byte("app.Start")) {
			opts.warnings = append(opts.warnings, "app.Start has been removed in v0.11.0. Use app.Serve Instead. [main.go]")
		}

		err = filepath.Walk(filepath.Join(opts.App.Root, "actions"), func(path string, info os.FileInfo, _ error) error {
			if info.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".go" {
				return nil
			}

			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if bytes.Contains(b, []byte("Websocket()")) {
				opts.warnings = append(opts.warnings,
					fmt.Sprintf("buffalo.Context#Websocket has been deprecated in v0.11.0, and removed in v0.12.0. Use github.com/gorilla/websocket directly. [%s]", path),
				)
			}
			if bytes.Contains(b, []byte("meta.Name")) {
				opts.warnings = append(opts.warnings,
					fmt.Sprintf("meta.Name has been deprecated in v0.11.0, and removed in v0.12.0. Use github.com/markbates/inflect.Name directly. [%s]", path),
				)
			}
			if bytes.Contains(b, []byte("generators.Find(")) {
				opts.warnings = append(opts.warnings,
					fmt.Sprintf("generators.Find(string) has been deprecated in v0.11.0, and removed in v0.12.0. Use generators.FindByBox() instead. [%s]", path),
				)
			}
			// i18n middleware changes in v0.11.1
			if bytes.Contains(b, []byte("T.CookieName")) {
				b = bytes.Replace(b, []byte("T.CookieName"), []byte("T.LanguageExtractorOptions[\"CookieName\"]"), -1)
			}
			if bytes.Contains(b, []byte("T.SessionName")) {
				b = bytes.Replace(b, []byte("T.SessionName"), []byte("T.LanguageExtractorOptions[\"SessionName\"]"), -1)
			}
			if bytes.Contains(b, []byte("T.LanguageFinder=")) || bytes.Contains(b, []byte("T.LanguageFinder ")) {
				opts.warnings = append(opts.warnings,
					fmt.Sprintf("i18n.Translator#LanguageFinder has been deprecated in v0.11.1, and has been removed in v0.12.0. Use i18n.Translator#LanguageExtractors instead. [%s]", path),
				)
			}
			return os.WriteFile(path, b, 0o664)
		})

		return err
	}
}
