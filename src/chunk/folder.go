package chunk

import (
	"os"
	"path/filepath"
)

func Paths(directory string) ([]string, error) {
	paths := []string{}

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	return paths, err
}