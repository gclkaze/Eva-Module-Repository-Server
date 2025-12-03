package utils

import (
	"os"
	"path/filepath"
)

func FolderExists(folder string) bool {
	_, err := os.Stat(folder)
	return !os.IsNotExist(err)
}

func CreateFolder(folder string) error {
	err := os.Mkdir(folder, os.ModePerm)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
	}
	return err
}

func CleanFolder(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		err := os.RemoveAll(path) // Removes file or directory recursively
		if err != nil {
			return err
		}
	}
	return nil
}
