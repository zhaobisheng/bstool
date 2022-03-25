package BCurl

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func getUrlRespHtml(strUrl, rtype string, postDict map[string]string) string {
	var respHtml string = ""
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}, // 使用环境变量的代理
		Proxy: http.ProxyFromEnvironment,
	}

	/*proxyUrl, err := url.Parse("http://127.0.0.1:8080")
	if err == nil { // 使用传入代理
		tr.Proxy = http.ProxyURL(proxyUrl)
		fmt.Println("Usage-proxyUrl:", proxyUrl)
	}*/

	httpClient := &http.Client{
		Transport: tr,
		//Transport:nil,
		//CheckRedirect: nil,
		//Jar:           CurCookieJar,
		CheckRedirect: MyCheckRedirect,
	}

	var httpReq *http.Request

	if nil == postDict {
		//httpReq, newReqErr = http.NewRequest("GET", strUrl, nil)
		httpReq, _ = http.NewRequest("GET", strUrl, nil)
		//httpReq.Header.Add("If-None-Match", `W/"wyzzy"`)
	} else {
		postValues := url.Values{}
		for postKey, PostValue := range postDict {
			postValues.Set(postKey, PostValue)
		}
		//fmt.Printf("postValues=%s\n", postValues)
		postDataStr := postValues.Encode()
		//fmt.Printf("postDataStr=%s\n", postDataStr)
		postDataBytes := []byte(postDataStr)
		//fmt.Printf("postDataBytes=%s\n", postDataBytes)
		postBytesReader := bytes.NewReader(postDataBytes)
		//httpReq, newReqErr = http.NewRequest("POST", strUrl, postBytesReader)
		httpReq, _ = http.NewRequest("POST", strUrl, postBytesReader)
		//httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	for key, val := range HeaderMap {
		httpReq.Header.Set(key, val)
	}
	if rtype == "first" {
		httpReq.Header.Set("Cookie", GetNewCookie())
	} else if rtype == "json" {
		httpReq.Header.Add("Content-Type", "application/json")
	}

	httpResp, err := httpClient.Do(httpReq)

	if err != nil {
		fmt.Printf("http get strUrl=%s response error=%s\n", strUrl, err.Error())
	}
	//fmt.Printf("httpResp.Header=%s\n", httpResp.Header)
	//fmt.Printf("httpResp.Status=%s\n", httpResp.Status)

	defer httpResp.Body.Close()

	body, errReadAll := ioutil.ReadAll(httpResp.Body)

	if errReadAll != nil {
		fmt.Printf("get response for strUrl=%s got error=%s\n", strUrl, errReadAll.Error())
	}

	//CurCookies = CurCookieJar.Cookies(httpReq.URL)

	//fmt.Println("请求session结果:", httpReq.URL, httpResp.Status, CurCookies)

	respHtml = string(body)
	return respHtml
}

func MyCheckRedirect(req *http.Request, via []*http.Request) error {
	//fmt.Printf("redirect:%v\n", req)
	if len(via) >= 10 {
		return errors.New("stopped after 10 redirects")
	}
	return nil
}
