package util

import (
	"os"
	"path/filepath"
)

// FindFiles is a convenience for walk a directory tree for a particular file
// name.
func FindFiles(root, name string) (paths []string, err error) {
	paths = []string{}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Base(path) == name {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return paths, nil
}
