package AliSMS

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const AccessKeyID = "xxxx"
const AccessKeySecret = "xxxxx"
const TemplateCode = "xxxxx"
const SignName = "搬砖人Bison"

func UrlEncode(in string) string {
	r := strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~")
	return r.Replace(url.QueryEscape(in))
}

func Sign(stringToSign string) string {
	h := hmac.New(sha1.New, []byte(fmt.Sprintf("%s&", AccessKeySecret)))
	h.Write([]byte(stringToSign))
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

func SortedString(data map[string]string) string {
	var sortQueryString string
	for _, v := range Keys(data) {
		sortQueryString = fmt.Sprintf("%s&%s=%s", sortQueryString, v, UrlEncode(data[v]))
	}
	return sortQueryString
}

func AliyunSendSMS(PhoneNumber, code string) ([]byte, error) {
	rand.Seed(time.Now().UnixNano())
	data := map[string]string{
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   strconv.FormatFloat(rand.Float64(), 'f', 6, 64),
		"AccessKeyId":      AccessKeyID,
		"SignatureVersion": "1.0",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Format":           "JSON",
		"RegionId":         "cn-hangzhou",
		"Action":           "SendSms",
		"Version":          "2017-05-25",
		"PhoneNumbers":     PhoneNumber,
		"SignName":         SignName,
		"TemplateCode":     TemplateCode,
		"TemplateParam":    `{"code":"` + code + `"}`,
	}

	sortQueryString := SortedString(data)
	stringToSign := "GET&%2F&" + UrlEncode(sortQueryString[1:])
	url := fmt.Sprintf("https://dysmsapi.aliyuncs.com/?Signature=%s%s", Sign(stringToSign), sortQueryString)

	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
