package utils

import (
	"os"
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
