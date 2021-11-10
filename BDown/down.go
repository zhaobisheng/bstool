package BDown

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func FileIsExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func Download(ur, dir, fn, proxyAddr, hash, hashType string, to time.Duration) error {
	if !FileIsExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println("mkdir "+dir+" err:", err.Error())
			return errors.New("mkdir " + dir + " err:" + err.Error())
		}
	}
	dfn := dir + "/" + fn
	var file *os.File
	var size int64
	if FileIsExist(dfn) {
		fi, err := os.OpenFile(dfn, os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println(fn, "open err:", err)
			return errors.New(fn + "open err:" + err.Error())
		}
		stat, _ := fi.Stat()
		size = stat.Size()
		sk, err := fi.Seek(size, 0)
		if err != nil {
			fmt.Println(fn, "seek err:", err)
			_ = fi.Close()
			return errors.New(fn + "seek err:" + err.Error())
		}
		if sk != size {
			fmt.Printf("%s seek length not equal file size,"+
				"seek=%d,size=%d\n", fn, sk, size)
			_ = fi.Close()
			return errors.New("seek length not equal file size")
		}
		file = fi
	} else {
		create, err := os.Create(dfn)
		if err != nil {
			fmt.Println(fn, "create err:", err)
			return errors.New(fn + "create err:" + err.Error())
		}
		file = create
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}, // 使用环境变量的代理
		Proxy: http.ProxyFromEnvironment,
	}
	if proxyAddr != "" {
		proxyUrl, err := url.Parse(proxyAddr)
		if err == nil { // 使用传入代理
			tr.Proxy = http.ProxyURL(proxyUrl)
			fmt.Println("Usage-proxy:", proxyUrl)
		}
	}
	client := &http.Client{Transport: tr}
	client.Timeout = to
	request := http.Request{}
	request.Method = http.MethodGet
	if size != 0 {
		header := http.Header{}
		header.Set("Range", "bytes="+strconv.FormatInt(size, 10)+"-")
		request.Header = header
	}
	parse, err := url.Parse(ur)
	if err != nil {
		fmt.Println(ur, "url err:", err)
		return errors.New(ur + "url err:" + err.Error())
	}
	request.URL = parse
	get, err := client.Do(&request)
	if err != nil {
		fmt.Println(ur, "get err:", err)
		return errors.New(ur + "get err:" + err.Error())
	}
	defer func() {
		err := get.Body.Close()
		if err != nil {
			fmt.Println(fn, "body close:", err.Error())
			//return errors.New(fn + "body err:" + err.Error())
		}
		err = file.Close()
		if err != nil {
			fmt.Println(fn, "file close:", err.Error())
			//return errors.New(fn + "file err:" + err.Error())
		}
	}()
	if get.ContentLength == 0 {
		fmt.Println(fn, "already downloaded")
		return nil
	}
	body := get.Body
	written, err := io.Copy(file, body)
	if err != nil {
		return errors.New("Copy-error:" + err.Error())
	}
	fmt.Println("written:", written)
	downSucc := true
	if hash != "" {
		var hashSum string
		var err error
		if strings.EqualFold(hashType, "md5") {
			hashSum, err = MD5sum(dfn)
		} else if strings.EqualFold(hashType, "sha1") {
			hashSum, err = Sha1sum(dfn)
		} else if strings.EqualFold(hashType, "sha256") {
			hashSum, err = Sha256sum(dfn)
		}
		if err == nil {
			if strings.EqualFold(hash, hashSum) {
				fmt.Println("Hash Verified!")
			} else {
				downSucc = false
				fmt.Println("Hash Error!")
			}
		}
	}
	if downSucc {
		fmt.Println(fn, "download success")
		return nil
	}
	return errors.New("Hash Error!")
}
