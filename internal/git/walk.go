package git

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// WalkDirs walks through directories to find .git folders and runs the provided callback on
// each.
func WalkDirs(dirs []string, callback func(path string) error) error {
	for _, dir := range dirs {
		err := filepath.WalkDir(
			dir,
			func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					if strings.Contains(err.Error(), "permission denied") {
						return nil
					}
					return err
				}

				if d.IsDir() && d.Name() == ".git" {
					if err := callback(path); err != nil {
						return err
					}
					return filepath.SkipDir
				}

				return nil
			},
		)
		if err != nil {
			return fmt.Errorf("Walking directory failed: %w", err)
		}
	}
	return nil
}
