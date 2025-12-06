package utils

import (
	"io/fs"
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

func ComputeFolderSizeBytes(root string) (int64, error) {
	var total int64

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Propagate errors (e.g., permission denied)
			return err
		}
		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Get FileInfo to access size
		info, err := d.Info()
		if err != nil {
			// If a file disappears or is unreadable, decide whether to ignore or fail.
			// Here we choose to fail.
			return err
		}

		// Only count regular files
		if info.Mode().IsRegular() {
			total += info.Size()
		}
		return nil
	})

	return total, err
}
