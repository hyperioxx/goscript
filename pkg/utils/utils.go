package utils

import (
	"os"
	"path/filepath"
)

func GetWorkingDirectory() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return currentDir, nil
}

func CheckFileExistsInDir(dir, fileName string) (string, bool) {
	filePath := filepath.Join(dir, fileName)
	_, err := os.Stat(filePath)
	return filePath, !os.IsNotExist(err)
}
