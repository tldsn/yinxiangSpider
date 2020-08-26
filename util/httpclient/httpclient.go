package httpclient

import (
	"compress/gzip"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	netUrl "net/url"
	"strings"
	"time"
)

//HPostFormHTML headers
func HPostFormHTML() map[string]string {
	ret := make(map[string]string)
	ret["Accept"] = "text/html, application/xhtml+xml, */*"
	ret["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; rv:12.0) Gecko/20120403211507 Firefox/12.0"
	ret["Connection"] = "Keep-Alive"
	ret["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	ret["Accept-Encoding"] = "gzip, deflate, br"
	ret["Accept-Language"] = "zh-CN,zh;q=0.8"
	return ret
}

//HPostFormJSON headers
func HPostFormJSON() map[string]string {
	ret := make(map[string]string)
	ret["Accept"] = "application/json, text/javascript, */*;"
	ret["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; rv:12.0) Gecko/20120403211507 Firefox/12.0"
	ret["Connection"] = "Keep-Alive"
	ret["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	ret["Accept-Encoding"] = "gzip, deflate, br"
	ret["Accept-Language"] = "zh-CN,zh;q=0.8"
	return ret
}

//HPostJSONJSON headers
func HPostJSONJSON() map[string]string {
	ret := make(map[string]string)
	ret["Accept"] = "application/json, */*;"
	ret["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:68.0) Gecko/20100101 Firefox/68.0"
	ret["Connection"] = "Keep-Alive"
	ret["Content-Type"] = "application/json;charset=UTF-8"
	ret["Accept-Encoding"] = "gzip, deflate"
	return ret
}

//HGetHTML headers
func HGetHTML() map[string]string {
	ret := make(map[string]string)
	ret["Accept"] = "text/html, application/xhtml+xml, */*"
	ret["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:68.0) Gecko/20100101 Firefox/68.0"
	ret["Connection"] = "Keep-Alive"
	ret["Accept-Encoding"] = "gzip, deflate"
	return ret
}

//HGetJSON headers
func HGetJSON() map[string]string {
	ret := make(map[string]string)
	ret["Accept"] = "application/json, text/plain, */*"
	ret["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:68.0) Gecko/20100101 Firefox/68.0"
	ret["Connection"] = "Keep-Alive"
	ret["Accept-Encoding"] = "gzip, deflate"
	return ret
}

//HMPostJSONJSON headers
func HMPostJSONJSON() map[string]string {
	ret := make(map[string]string)
	ret["Accept"] = "application/json, */*;"
	ret["User-Agent"] = "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko)"
	ret["Connection"] = "Keep-Alive"
	ret["Content-Type"] = "application/json;charset=UTF-8"
	ret["Accept-Encoding"] = "gzip, deflate"
	return ret
}

//HMPostFormJSON headers
func HMPostFormJSON() map[string]string {
	ret := make(map[string]string)
	ret["Accept"] = "application/json, */*;"
	ret["User-Agent"] = "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko)"
	ret["Connection"] = "Keep-Alive"
	ret["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	ret["Accept-Encoding"] = "gzip, deflate"
	return ret
}

//HMGetJSON headers
func HMGetJSON() map[string]string {
	ret := make(map[string]string)
	ret["Accept"] = "application/json;*/*;"
	ret["User-Agent"] = "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko)"
	ret["Connection"] = "Keep-Alive"
	ret["Accept-Encoding"] = "gzip, deflate"
	return ret
}

//Get ...
func Get(url, proxy string, headers map[string]string, timeout int) (map[string]string, error) {
	return DoRequest("GET", url, proxy, nil, headers, timeout)
}

//Post ...
func Post(url, post, proxy string, headers map[string]string, timeout int) (map[string]string, error) {
	// return DoRequest("POST", url, proxy, strings.NewReader(post), headers, timeout)
	return DoRequest2("", "POST", url, post, headers, timeout, proxy, false)

}

func DoRequest2(httpType string, method string, url string, param string, headers map[string]string, timeout int, proxyUrl string, redirect bool) (map[string]string, error) {

	var transport http.RoundTripper

	transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	if proxyUrl != "" {
		p := ""
		if strings.Contains(proxyUrl, "|") {
			arr := strings.Split(proxyUrl, "|")
			usr := arr[0]
			pwd := arr[1]
			proxyUrl = arr[2]
			u, err := netUrl.Parse(url)
			if err != nil {
				return nil, err
			}
			p = u.Scheme + "://" + usr + ":" + pwd + "@" + proxyUrl
		}
		transport.(*http.Transport).Proxy = func(i *http.Request) (*netUrl.URL, error) {
			if p != "" {
				return netUrl.Parse(p)
			} else {
				return &netUrl.URL{
					Host: proxyUrl,
				}, nil
			}
		}
	}

	// if httpType == HTTP2 {
	// 	err := http2.ConfigureTransport(transport.(*http.Transport))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	client := http.DefaultClient
	client.Transport = transport
	// 指定CheckRedirect函数，不进行重定向操作
	if !redirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client.Jar = jar
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Millisecond
	}
	req, err := http.NewRequest(method, url, strings.NewReader(param))
	if err != nil {
		return nil, err
	}

	if nil != headers && len(headers) != 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	rsp, err := client.Do(req)
	if nil != err {
		return nil, err
	}

	return parseResp(rsp)
}

//DoRequest ...
func DoRequest(method, urlstr, proxy string, body io.Reader, headers map[string]string, timeout int) (map[string]string, error) {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true, Renegotiation: tls.RenegotiateOnceAsClient}}
	if proxy != "" {
		proxyURL, err := netUrl.Parse(proxy)
		if err != nil {
			log.Println(err)
		}
		tr.Proxy = http.ProxyURL(proxyURL)
	}
	client := &http.Client{Transport: tr}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Millisecond
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client.Jar = jar
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest(method, urlstr, body)
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return parseResp(resp)
}

//ProcessCookie ...
func ProcessCookie(cookie string) string {
	arr := strings.Split(cookie, ";")
	mapTemp := make(map[string]string)
	for _, val := range arr {
		index := strings.Index(val, "=")
		if index == -1 {
			continue
		}
		mapTemp[val[:index]] = val
	}
	ret := ""
	for _, val := range mapTemp {
		ret += val + ";"
	}
	return ret
}

func parseResp(resp *http.Response) (map[string]string, error) {
	ret := make(map[string]string)
	cookies := ""
	for _, val := range resp.Cookies() {
		cookies += val.Name + "=" + val.Value + ";"
	}
	cookies = ProcessCookie(cookies)
	isGzip := false
	for key, values := range resp.Header {
		if strings.EqualFold(key, "Content-Encoding") {
			if values[0] == "gzip" {
				isGzip = true
			}
		} else if !strings.EqualFold(key, "Set-Cookie") {
			ret[key] = values[0]
		}
	}
	ret["cookie"] = cookies

	defer resp.Body.Close()
	if isGzip {
		gzipR, err := gzip.NewReader(resp.Body)
		body, err := ioutil.ReadAll(gzipR)
		if err != nil {
			return nil, err
		}

		ret["body"] = string(body[:])
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		ret["body"] = string(body[:])
	}

	return ret, nil
}
