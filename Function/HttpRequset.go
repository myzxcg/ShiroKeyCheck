package Function

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"ShiroKeyCheck/GlobalVar"
)

func HttpRequset(RememberMe string) bool {
	//设置跳过https证书验证，超时和代理
	var tr *http.Transport
	if GlobalVar.HttpProxy != "" {
		uri, _ := url.Parse(GlobalVar.HttpProxy)
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(uri),
		}
	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	client := &http.Client{
		Timeout:   time.Duration(GlobalVar.Timeout) * time.Second,
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse //不允许跳转
		}}
	req, err := http.NewRequest(strings.ToUpper(GlobalVar.Method), GlobalVar.Url, strings.NewReader(GlobalVar.PostContent))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//设置请求头
	if strings.ToUpper(GlobalVar.Method) == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("User-Agent", GlobalVar.UserAgent)
	req.Header.Set("Cookie", "rememberMe="+RememberMe)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	//判断rememberMe=deleteMe;是否在响应头中
	var SetCookieAll string
	for i := range resp.Header["Set-Cookie"] {
		SetCookieAll += resp.Header["Set-Cookie"][i]
	}
	if !strings.Contains(SetCookieAll, "rememberMe=deleteMe;") {
		return true //内容中不包含rememberMe
	} else {
		return false
	}
}
