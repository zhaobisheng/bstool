package BDown

import (
	"bufio"
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

func Download(ur, dir, fn, proxyAddr, md5 string, to time.Duration, onErr func(err error)) {
	if !FileIsExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println("mkdir "+dir+" err:", err.Error())
			onErr(err)
			return
		}
	}
	dfn := dir + "/" + fn
	var file *os.File
	var size int64
	if FileIsExist(dfn) {
		fi, err := os.OpenFile(dfn, os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println(fn, "open err:", err)
			onErr(err)
			return
		}
		stat, _ := fi.Stat()
		size = stat.Size()
		sk, err := fi.Seek(size, 0)
		if err != nil {
			fmt.Println(fn, "seek err:", err)
			_ = fi.Close()
			onErr(err)
			return
		}
		if sk != size {
			fmt.Printf("%s seek length not equal file size,"+
				"seek=%d,size=%d\n", fn, sk, size)
			_ = fi.Close()
			onErr(errors.New("seek length not equal file size"))
			return
		}
		file = fi
	} else {
		create, err := os.Create(dfn)
		if err != nil {
			fmt.Println(fn, "create err:", err)
			onErr(err)
			return
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
		onErr(err)
		return
	}
	request.URL = parse
	get, err := client.Do(&request)
	if err != nil {
		fmt.Println(ur, "get err:", err)
		onErr(err)
		return
	}
	defer func() {
		err := get.Body.Close()
		if err != nil {
			fmt.Println(fn, "body close:", err.Error())
			onErr(err)
		}
		err = file.Close()
		if err != nil {
			fmt.Println(fn, "file close:", err.Error())
			onErr(err)
		}
	}()
	if get.ContentLength == 0 {
		fmt.Println(fn, "already downloaded")
		return
	}
	body := get.Body
	writer := bufio.NewWriter(file)
	bs := make([]byte, 10*1024*1024) //每次读取的最大字节数，不可为0
	for {
		var read int
		read, err = body.Read(bs)
		if err != nil {
			if err != io.EOF {
				fmt.Println(fn, "read err:"+err.Error())
				onErr(err)
			} else {
				err = nil
			}
			break
		}
		_, err = writer.Write(bs[:read])
		if err != nil {
			fmt.Println(fn, "write err:"+err.Error())
			onErr(err)
			break
		}
	}
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println(fn, "writer flush:", err.Error())
		onErr(err)
		return
	}

	fmt.Println(fn, "download success")

	if md5 != "" {
		md5sum, err := MD5sum(dfn)
		if err == nil {
			if strings.EqualFold(md5, md5sum) {
				fmt.Println("md5 Verified!")
			}
		}
	}
}
