package FileUtils

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Write2file(filename, str string) {
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	defer fp.Close()
	if err == nil {
		fp.WriteString(str)
	}
}

func ReadLineFile(filename string, showProcess bool, callback func(string)) {
	fi, err := os.Open(filename)
	if err != nil {
		log.Printf("read-ipfile-error: %s\n", err)
		return
	}
	defer func() {
		fi.Close()
	}()
	br := bufio.NewReader(fi)
	Num := 0
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			log.Println("文件读取完成,共有", Num, "条")
			break
		}
		Num++
		lineStr := string(a)
		callback(lineStr)
		if showProcess {
			if Num%1000 == 0 {
				log.Println("已读取", Num, "条")
			}
		}
	}
}

func ReadHeaderFromFile(filename string) map[string]string {
	tempHeaderMap := make(map[string]string)
	fi, err := os.Open(filename)
	if err != nil {
		log.Printf("read-header-error: %s\n", err)
		return tempHeaderMap
	}
	defer func() {
		fi.Close()
	}()
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		headerStr := string(a)
		rs := strings.Split(string(headerStr), ": ")
		for index := 0; index <= len(rs)/2; index += 2 {
			headerKey := rs[index]
			tempHeaderMap[headerKey] = rs[index+1]
		}
	}
	return tempHeaderMap
}

func FileList(path string) []os.FileInfo {
	dir_list, e := ioutil.ReadDir(path)
	if e != nil {
		return nil
	}
	return dir_list
}

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
