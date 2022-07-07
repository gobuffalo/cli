package plugin

import (
	"os"
	"path/filepath"
)

// WorkDirValidator checks whether the current working directory
// is ok for a plugin to run. This is specially useful for commands
// that are intended to run within the Buffalo project.
type WorkDirValidator interface {
	// ValidateWorkDir receives the Work directory and returns
	// if its valid as well as any error it could have found.
	ValidateWorkDir(wd string) (bool, error)
}

func ValidateBuffaloRoot(wd string) (bool, error) {
	_, err := os.Stat(filepath.Join(wd, ".buffalo.dev.yml"))
	if err != nil {
		return false, nil
	}

	return true, nil
}
