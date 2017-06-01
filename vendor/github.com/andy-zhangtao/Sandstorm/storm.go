package Sandstorm

import (
	"errors"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

var cookies []*http.Cookie

var (
	debug = true
)

//HTTPError 返回错误信息
func HTTPError(w http.ResponseWriter, err string, status int) {
	http.Error(w, err, status)
}

//HTTPSuccess 返回成功信息
func HTTPSuccess(w http.ResponseWriter, describe string) {
	w.Write([]byte(describe))
}

func HTTPReDirect(w http.ResponseWriter, r *http.Request, header map[string]string) {
	// w.Header().Set("Location", header["Location"])
	http.Redirect(w, r, header["Location"], http.StatusMovedPermanently)
}

//Get Get方法
//schema: HTTP/HTTPS
//u: url
//args: GET方法查询参数
func Get(u string, args map[string]string, proxy string) (*http.Response, string, error) {
	return _send(http.MethodGet, u, args, "", proxy)
}

//Put Put方法
//schema: HTTP/HTTPS
//u: url
//args: PUT方法参数
func Put(u string, args1 map[string]string, args2 string, proxy string) (*http.Response, string, error) {
	return _send(http.MethodPut, u, args1, args2, proxy)
}

//Post Post
//schema: HTTP/HTTPS
//u: url
//args: Post方法参数
func Post(u string, args1 map[string]string, args2 string, proxy string) (*http.Response, string, error) {
	return _send(http.MethodPost, u, args1, args2, proxy)
}

//Delete Delete
//schema: HTTP/HTTPS
//u: url
//args: Delete方法参数
func Delete(u string, args1 string, proxy string) (*http.Response, string, error) {
	return _send(http.MethodDelete, u, nil, args1, proxy)
}

func _send(method string, u string, args map[string]string, json string, proxy string) (*http.Response, string, error) {
	request := gorequest.New()
	request.SetDebug(debug)
	if proxy != "" {
		request.Proxy(proxy)
	}
	request.Set("User-Agent", "go_client_http_0.1")
	if len(cookies) > 0 {
		request.AddCookies(cookies)
	}

	var resp *http.Response
	body := ""
	var errs []error

	switch method {
	case http.MethodPut:
		resp, body, errs = request.Put(u).Send(json).End()
		break
	case http.MethodPost:
		resp, body, errs = request.Post(u).Send(json).End()
		break
	case http.MethodDelete:
		resp, body, errs = request.Delete(u).Send(json).End()
		break
	case http.MethodGet:
		superAgent := request.Get(u)
		for k := range args {
			query := k + "=" + args[k]
			superAgent = superAgent.Query(query)
		}
		resp, body, errs = superAgent.End()
		break
	}
	if len(errs) > 0 {
		err := ""
		for _, e := range errs {
			err = err + e.Error()
		}
		return nil, "", errors.New(err)
	}

	return resp, body, nil

}

// SetCookie 设置cookie
func SetCookie(c map[string]string) {
	if cookies == nil {
		cookies = make([]*http.Cookie, len(c))
	}
	i := 0
	for k := range c {
		ck := http.Cookie{
			Name:  k,
			Value: c[k],
		}
		cookies[i] = &ck
		i++
	}

}

// DisDebug 关闭debug输出
func DisDebug() {
	debug = false
}

// EnDebug 启动debug输出
func EnDebug() {
	debug = true
}
