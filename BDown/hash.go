// Package md5 computes MD5 checksum for large files

package BDown

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

const bufferSize = 65536

// MD5sum returns MD5 checksum of filename
func MD5sum(filename string) (string, error) {
	if info, err := os.Stat(filename); err != nil {
		return "", err
	} else if info.IsDir() {
		return "", nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	for buf, reader := make([]byte, bufferSize), bufio.NewReader(file); ; {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		hash.Write(buf[:n])
	}
	checksum := fmt.Sprintf("%x", hash.Sum(nil))
	return checksum, nil
}

/*
func Md5Str(data string) string {
	bytes := []byte(data)
	hash := md5.New()
	hash.Write(bytes)
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}
*/
func Sha1sum(filename string) (string, error) {
	if info, err := os.Stat(filename); err != nil {
		return "", err
	} else if info.IsDir() {
		return "", nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha1.New()
	for buf, reader := make([]byte, bufferSize), bufio.NewReader(file); ; {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		hash.Write(buf[:n])
	}
	checksum := fmt.Sprintf("%x", hash.Sum(nil))
	return checksum, nil
}

func Sha256sum(filename string) (string, error) {
	if info, err := os.Stat(filename); err != nil {
		return "", err
	} else if info.IsDir() {
		return "", nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	for buf, reader := make([]byte, bufferSize), bufio.NewReader(file); ; {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		hash.Write(buf[:n])
	}
	checksum := fmt.Sprintf("%x", hash.Sum(nil))
	return checksum, nil
}

/*


package main

import (

"crypto/md5"

"fmt"

"io"

"math"

"os"

)

const filechunk = 8192 // we settle for 8KB

func main() {
file, err := os.Open("utf8.txt")

if err != nil {
panic(err.Error())

}

defer file.Close()

// calculate the file size

info, _ := file.Stat()

filesize := info.Size()

blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))

hash := md5.New()

for i := uint64(0); i < blocks; i++ {
blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))

buf := make([]byte, blocksize)

file.Read(buf)

io.WriteString(hash, string(buf)) // append into the hash

}

fmt.Printf("%s checksum is %x\n", file.Name(), hash.Sum(nil))

}
*/
