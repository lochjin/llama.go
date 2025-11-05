package common

import (
	"os"
	"path/filepath"
	"strings"
)

func ReadFile(path string) ([]byte, error) {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	ba, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ba, nil
}

func IsFilePath(str string) bool {
	if filepath.IsAbs(str) {
		return true
	}
	if strings.ContainsRune(str, filepath.Separator) {
		return true
	}
	return false
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func SaveOutputToFile(outFilePath string, content string) error {
	outFile, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		outFile.Close()
	}()
	_, err = outFile.WriteString(content)
	return err
}
