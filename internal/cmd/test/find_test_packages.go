package test

import "golang.org/x/tools/go/packages"

// findTestPackages in the current directory using the go/packages API.
func findTestPackages(givenArgs []string) ([]string, error) {
	// If there are args, then assume these are the packages to test.
	if len(givenArgs) > 0 {
		return givenArgs, nil
	}

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
		xx = append(xx, pkg.ID)
	}

	return xx, nil
}
