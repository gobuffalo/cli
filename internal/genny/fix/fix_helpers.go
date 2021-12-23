package fix

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
)

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

		// skip directories that are not part of the buffalo source tree
		for _, prefix := range []string{"vendor", "node_modules", ".git"} {
			if strings.HasPrefix(f.Name(), prefix) {
				return nil
			}
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
