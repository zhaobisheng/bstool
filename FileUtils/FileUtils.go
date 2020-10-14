package FileUtils

import (
	"io/ioutil"
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

func CreateDirIfNotExist(dirName string) bool {
	if !PathExists(dirName) {
		err := os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			return false
		}
	}
	return true
}

func CopyFile(srcFile, newFile string) error {
	fileData, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(newFile, fileData, os.ModePerm)
	return err
}