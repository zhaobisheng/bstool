package Curl

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

func HandleParam(tempMap map[string]string) []byte {
	urlParam := ""
	for key, val := range tempMap {
		urlParam += key + "=" + val + "&"
	}
	urlParam = TrimFlag(urlParam, "&")
	return []byte(urlParam)
}

//"127.0.0.1:10808"
func HttpSock5(turl, method, addr string, data []byte, header map[string]string) (resp *http.Response, err error) {
	dialer, err := proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxyUrl:", err)
		os.Exit(1)
	}
	httpTransport := &http.Transport{TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	}}
	client := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial
	request, err := http.NewRequest(method, turl, bytes.NewReader(data))
	if err != nil {
		fmt.Println("err:", err)
	}
	if header != nil {
		for key, val := range header {
			request.Header.Set(key, val)
		}
	} else {
		request.Header.Set("Content-Type", "application/json")
	}
	resp, err = client.Do(request)
	return resp, err
}

func HttpRequest(turl, method, proxyUrl string, data []byte, header map[string]string) (resp *http.Response, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}, // 使用环境变量的代理
		Proxy: http.ProxyFromEnvironment,
	}
	if proxyUrl != "" {
		proxyUrl, err := url.Parse(proxyUrl)
		if err == nil { // 使用传入代理
			tr.Proxy = http.ProxyURL(proxyUrl)
			fmt.Println("Usage-proxyUrl:", proxyUrl)
		}
	}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest(method, turl, bytes.NewReader(data))
	if err != nil {
		fmt.Println("err:", err)
	}
	if header != nil {
		for key, val := range header {
			if strings.EqualFold(key, "host") {
				request.Host = val
			}
			request.Header.Set(key, val)
		}
	} else {
		request.Header.Set("Content-Type", "application/json")
	}
	resp, err = client.Do(request)
	return resp, err
}

func RequestResult(url, method, proxyUrl string, data []byte, header map[string]string, sock5 bool) (string, error) {
	var response *http.Response
	var err error
	if sock5 {
		response, err = HttpSock5(url, method, proxyUrl, data, header)
	} else {
		response, err = HttpRequest(url, method, proxyUrl, data, header)
	}
	if err != nil {
		return "", err
	}
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func DownloadFile(url, method, downPath, proxyUrl string, data []byte, header map[string]string) error {
	response, err := HttpRequest(url, method, proxyUrl, data, header)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(downPath+"/"+GetFilename(url), buf, os.ModePerm)
	return err
}

func GetFileData(url, method, proxyUrl string, data []byte, header map[string]string) []byte {
	response, err := HttpRequest(url, method, proxyUrl, data, header)
	if err != nil {
		return []byte("")
	}
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte("")
	}
	return buf
}

func GetFilename(targetUrl string) string {
	index := strings.LastIndex(targetUrl, "/")
	if index > 0 {
		return targetUrl[index+1:]
	}
	return fmt.Sprintf("%v", time.Now().Unix())
}

var userAgent = [...]string{"Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, 360SE)",
	"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
	"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
	"Opera/9.80 (Windows NT 6.1, U, en) Presto/2.8.131 Version/11.11",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, TencentTraveler 4.0)",
	"Mozilla/5.0 (Windows, U, Windows NT 6.1, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
	"Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Linux, U, Android 3.0, en-us, Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13",
	"Mozilla/5.0 (iPad, U, CPU OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, Trident/4.0, SE 2.X MetaSr 1.0, SE 2.X MetaSr 1.0, .NET CLR 2.0.50727, SE 2.X MetaSr 1.0)",
	"Mozilla/5.0 (iPhone, U, CPU iPhone OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
	"MQQBrowser/26 Mozilla/5.0 (Linux, U, Android 2.3.7, zh-cn, MB200 Build/GRJ22, CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandomUserAgent() string {
	return userAgent[r.Intn(len(userAgent))]
}

func TrimFlag(str, symbol string) string {
	newStr := str
	if len(newStr) > len(symbol) {
		flagLen := len(symbol)
		if str[len(str)-flagLen:] == symbol {
			newStr = str[:len(str)-flagLen]
		}
	}
	return newStr
}
