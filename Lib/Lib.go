package Lib

import (
	"fmt"

	"math/rand"
	"net"
	"net/http"

	"os"

	"time"

	"github.com/axgle/mahonia"
)

func CheckResponse(code int) bool {
	if code == 403 {
		fmt.Println("403-没权限访问")
		return false
	}
	if code == 503 {
		fmt.Println("503-IP已被禁止访问")
		return false
	}
	if code == 200 {
		return true
	}
	if code == 404 {
		fmt.Println("404-找不到页面")
		return false
	}
	return false
}

/**
* 返回response
 */
func GetResponse(url string) (*http.Response, error) {
	fmt.Println(url)
	client := &http.Client{}
	client.Timeout = 60 * time.Second
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 60 * time.Second,
		}).Dial,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   60 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
	}
	client.Transport = transport
	request, _ := http.NewRequest("GET", url, nil)
	User_Agent := GetRandomUserAgent()
	request.Header.Set("User-Agent", User_Agent)
	response, err := client.Do(request)
	/*for {
		if response != nil {
			if CheckResponse(response.StatusCode) {
				break
			} else {
				tc := time.After(5 * time.Second)
				fmt.Println("Waiting 5 second!")
				<-tc
				User_Agent = GetRandomUserAgent()
				request.Header.Set("User-Agent", User_Agent)
				response, _ = client.Do(request)
			}
		} else {
			tc := time.After(3 * time.Second)
			fmt.Println("Waiting 3 second!")
			<-tc
			User_Agent = GetRandomUserAgent()
			request.Header.Set("User-Agent", User_Agent)
			response, _ = client.Do(request)
		}
	}*/
	return response, err
}

var userAgent = [...]string{"Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, 360SE)",
	"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
	"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
	"Opera/9.80 (Windows NT 6.1, U, en) Presto/2.8.131 Version/11.11",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, TencentTraveler 4.0)",
	"Mozilla/5.0 (Windows, U, Windows NT 6.1, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/19.0.963.56 Safari/535.11",
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

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
