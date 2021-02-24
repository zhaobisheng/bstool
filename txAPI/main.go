// txDomain project main.go
package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

const apiUrl = "cns.api.qcloud.com/v2/index.php"

func TXRecordCreate(AccessKeyID, AccessKeySecret, domain string) string {
	rand.Seed(time.Now().UnixNano())
	data := map[string]string{
		"Nonce":     fmt.Sprintf("%d", rand.Int31()),
		"SecretId":  AccessKeyID,
		"Timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		"Action":    "RecordList",
		"domain":    domain, //"funwan.cn",
		"offset":    "0",
		"length":    "20",
		"Region":    "ap-guangzhou",
	}

	sortQueryString := SortedString(data)
	stringToSign := "POST" + apiUrl + "?" + sortQueryString[1:]
	//fmt.Println(stringToSign)
	Signature := Sign(stringToSign, AccessKeySecret)
	//url := fmt.Sprintf("https://%sSignature=%s%s", apiUrl, Signature, sortQueryString)
	turl := fmt.Sprintf("https://%s", apiUrl)
	//fmt.Println(turl)
	data["Signature"] = Signature
	sortQueryString = SortedString(data)
	param := sortQueryString[1:]
	//fmt.Println(param)
	/*iurlproxy, _ := url.Parse("http://127.0.0.1:8080")
	tr := &http.Transport{
		TLSClentConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(urlproxy),
	}
	c := &http.Client{Transport: tr}*/
	r, err := http.Post(turl, "application/x-www-form-urlencoded", strings.NewReader(param))
	if err != nil {
		fmt.Println("post-error:", err)
		return ""
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	return string(b)
}

func SortedString(data map[string]string) string {
	var sortQueryString string
	for _, v := range Keys(data) {
		sortQueryString = fmt.Sprintf("%s&%s=%s", sortQueryString, v, data[v])
		//sortQueryString = fmt.Sprintf("%s&%s=%s", sortQueryString, v, UrlEncode(data[v]))
	}
	return sortQueryString
}

func UrlEncode(in string) string {
	r := strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~")
	return r.Replace(url.QueryEscape(in))
}

func Sign(stringToSign, AccessKeySecret string) string {
	h := hmac.New(sha1.New, []byte(AccessKeySecret))
	h.Write([]byte(stringToSign))
	fmt.Println(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	return UrlEncode(base64.StdEncoding.EncodeToString(h.Sum(nil)))
}

func Keys(data map[string]string) []string {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
