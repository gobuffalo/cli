package test

import (
	"golang.org/x/tools/go/packages"
)

// findPackages in the current directory using the x/tools/go/packages API
func findPackages() ([]string, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName,
		Dir:  ".",
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return []string{}, err
	}

	if packages.PrintErrors(pkgs) > 0 {
		return []string{}, err
	}

	xx := []string{}
	for _, pkg := range pkgs {
		xx = append(
			xx,

			// Trim the prefix of the application root module
			pkg.ID,
		)
	}

	return xx, nil
}
