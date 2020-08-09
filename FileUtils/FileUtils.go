package FileUtils

import (
	"os"
	"path/filepath"
)

func GetFileDir(path string) string {
	return filepath.Dir(path)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
