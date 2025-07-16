package builder

import (
	"io/fs"
	"path/filepath"
)

func ListFilesRecursively(root string) ([]string, error) {
	var result []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		result = append(result, path)
		return nil
	})
	return result, err
}