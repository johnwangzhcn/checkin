package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	gladosUrl      = "https://glados.rocks/api/user/checkin"
	pushToken      = "a5c767e0ffde4b13917f825e508448be"
	v2freeLoginUrl = "https://w1.v2free.top/auth/login"
	v2freeCheckUrl = "https://w1.v2free.net/user/checkin"
	v2freeUserUrl  = "https://w1.v2free.net/user"
)

type Account struct {
	ID       string
	Cookie   string
	Password string
	Req      Request
}

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ret     int    `json:"ret"`
	Msg     string `json:"msg"`
}

var gladosAccounts = []Account{
	{
		ID:     "xxx@88.com",
		Cookie: "cf_clearance=WGyq9l1rvt8tsFE2NvhFUWVSzw1TnPzmQb7ZidJXJ_w-1659352632-0-150; koa:sess=eyJ1c2VySWQiOjEwNTIxNCwiY29kZSI6IjczMDhZLUIyWjhOLU1KUEw5LU1IWUZCIiwiX2V4cGlyZSI6MTY4OTc4NTg0MjM4OSwiX21heEFnZSI6MjU5MjAwMDAwMDB9; koa:sess.sig=qJhcKdEAXuG_xXh7NrB3JLNr_mY; googtrans=/auto/zh-CN; googtrans=/auto/zh-CN",
		Req: Request{
			URL:     gladosUrl,
			Method:  http.MethodPost,
			Payload: strings.NewReader(`{"token":"glados.network"}`),
			Headers: map[string]string{"content-type": "application/json;charset=UTF-8"},
		},
	},
	{
		ID:     "orz3-8@protonmail.com",
		Cookie: "_ga=GA1.2.530142443.1639584545; koa:sess=eyJ1c2VySWQiOjE5ODYyNCwiX2V4cGlyZSI6MTY4Nzk0NDkxMTc1NiwiX21heEFnZSI6MjU5MjAwMDAwMDB9; koa:sess.sig=1U9ei2WJue4LSErjLexACMLrhb8; __stripe_mid=25344988-3253-404d-8986-0b2a009f5086805478; _gid=GA1.2.639938267.1673534053; _gat_gtag_UA_104464600_2=1",
		Req: Request{
			URL:     gladosUrl,
			Method:  http.MethodPost,
			Payload: strings.NewReader(`{"token":"glados.network"}`),
			Headers: map[string]string{"content-type": "application/json;charset=UTF-8"},
		},
	},
	{
		ID:     "xxx@qq.com",
		Cookie: "koa:sess=eyJ1c2VySWQiOjEwNTU1NSwiX2V4cGlyZSI6MTY4Nzk3NjAwOTk0MiwiX21heEFnZSI6MjU5MjAwMDAwMDB9; koa:sess.sig=X6vi4rnD7BMbqWX-eK-1214c7D4",
		Req: Request{
			URL:     gladosUrl,
			Method:  http.MethodPost,
			Payload: strings.NewReader(`{"token":"glados.network"}`),
			Headers: map[string]string{"content-type": "application/json;charset=UTF-8"},
		},
	},
	{
		ID:     "orz3-8@pm.me",
		Cookie: "_ga=GA1.2.1413178732.1673089766; koa:sess=eyJ1c2VySWQiOjIwODQ5NSwiX2V4cGlyZSI6MTY5OTAwOTg0NjA3OCwiX21heEFnZSI6MjU5MjAwMDAwMDB9; koa:sess.sig=vfLLiyOW5jF43iplgSy24hLUbsE; _gid=GA1.2.33136766.1673537618",
		Req: Request{
			URL:     gladosUrl,
			Method:  http.MethodPost,
			Payload: strings.NewReader(`{"token":"glados.network"}`),
			Headers: map[string]string{"content-type": "application/json;charset=UTF-8"},
		},
	},
	{
		ID:     "xxx@gmail.com",
		Cookie: "koa:sess=eyJ1c2VySWQiOjI1MzA2NCwiX2V4cGlyZSI6MTcwMTQyMDI1MTAxNCwiX21heEFnZSI6MjU5MjAwMDAwMDB9; koa:sess.sig=U2Iw4sObS6z6Lt4QV42jpZkkbnw",
		Req: Request{
			URL:     gladosUrl,
			Method:  http.MethodPost,
			Payload: strings.NewReader(`{"token":"glados.network"}`),
			Headers: map[string]string{"content-type": "application/json;charset=UTF-8"},
		},
	},
	{
		ID:     "xxx@hotmail.com",
		Cookie: "_ga=GA1.1.178142109.1676899058; koa:sess=eyJ1c2VySWQiOjI4MDE4MiwiX2V4cGlyZSI6MTcwMjgxOTEwNjcyNCwiX21heEFnZSI6MjU5MjAwMDAwMDB9; koa:sess.sig=llMHBwoxZeyhj93e_0Oryq1bIJE; _ga_CZFVKMNT9J=GS1.1.1676899057.1.1.1676899216.0.0.0",
		Req: Request{
			URL:     gladosUrl,
			Method:  http.MethodPost,
			Payload: strings.NewReader(`{"token":"glados.network"}`),
			Headers: map[string]string{"content-type": "application/json;charset=UTF-8"},
		},
	},
	{
		ID:     "orz3-8@proton.me",
		Cookie: "koa:sess=eyJ1c2VySWQiOjI0MzA4NywiX2V4cGlyZSI6MTcwMTYyMTI5ODk0NiwiX21heEFnZSI6MjU5MjAwMDAwMDB9; koa:sess.sig=o0dzSX-CHzmA0zYhFckJ8CL4S-g; __stripe_mid=ddb3f6d6-8bdd-4df2-9959-c3837ac320d1fe9cbe; _ga_CZFVKMNT9J=GS1.1.1678104474.1.1.1678105227.0.0.0; _ga=GA1.2.1597457881.1678104474; _gid=GA1.2.85789393.1678104475",
		Req: Request{
			URL:     gladosUrl,
			Method:  http.MethodPost,
			Payload: strings.NewReader(`{"token":"glados.network"}`),
			Headers: map[string]string{"content-type": "application/json;charset=UTF-8"},
		},
	},
	{
		ID:     "y@pm",
		Cookie: "koa:sess=eyJ1c2VySWQiOjI5MzgxMiwiX2V4cGlyZSI6MTcwNDQ3NDQ3NjIxMywiX21heEFnZSI6MjU5MjAwMDAwMDB9; koa:sess.sig=btyqVA_PF3SzLd6H4q8eqtywAYk",
		Req: Request{
			URL:     gladosUrl,
			Method:  http.MethodPost,
			Payload: strings.NewReader(`{"token":"glados.network"}`),
			Headers: map[string]string{"content-type": "application/json;charset=UTF-8"},
		},
	},
	{
		ID:     "johnwong@chinamail.com",
		Cookie: "_gid=GA1.2.7057315.1681145872; koa:sess=eyJ1c2VySWQiOjMyMzc2NCwiX2V4cGlyZSI6MTcwNzExODE3MjcwOSwiX21heEFnZSI6MjU5MjAwMDAwMDB9; koa:sess.sig=CE1lR7kOjtoGjd4pGSeddzGtZxk; _ga=GA1.1.802715368.1681145872; _ga_CZFVKMNT9J=GS1.1.1681198024.3.1.1681198273.0.0.0",
		Req: Request{
			URL:     gladosUrl,
			Method:  http.MethodPost,
			Payload: strings.NewReader(`{"token":"glados.network"}`),
			Headers: map[string]string{"content-type": "application/json;charset=UTF-8"},
		},
	},
}

var v2freeAccounts = []Account{
	{
		ID:       "john@free.com",
		Password: "2583585355",
		Req: Request{
			URL:    v2freeLoginUrl,
			Method: http.MethodPost,
			Headers: map[string]string{
				"content-type": "application/x-www-form-urlencoded; charset=UTF-8",
				"referer":      v2freeLoginUrl,
			},
		},
	},
	{
		ID:       "johnson@free.com",
		Password: "2583585355",
		Req: Request{
			URL:    v2freeLoginUrl,
			Method: http.MethodPost,
			Headers: map[string]string{
				"content-type": "application/x-www-form-urlencoded; charset=UTF-8",
				"referer":      v2freeLoginUrl,
			},
		},
	},
	{
		ID:       "orz3-8@protonmail.com",
		Password: "2583585355",
		Req: Request{
			URL:    v2freeLoginUrl,
			Method: http.MethodPost,
			Headers: map[string]string{
				"content-type": "application/x-www-form-urlencoded; charset=UTF-8",
				"referer":      v2freeLoginUrl,
			},
		},
	},
	{
		ID:       "orz3-8@proton.me",
		Password: "2583585355",
		Req: Request{
			URL:    v2freeLoginUrl,
			Method: http.MethodPost,
			Headers: map[string]string{
				"content-type": "application/x-www-form-urlencoded; charset=UTF-8",
				"referer":      v2freeLoginUrl,
			},
		},
	},
	{
		ID:       "john@wong.cn",
		Password: "johnwong",
		Req: Request{
			URL:    v2freeLoginUrl,
			Method: http.MethodPost,
			Headers: map[string]string{
				"content-type": "application/x-www-form-urlencoded; charset=UTF-8",
				"referer":      v2freeLoginUrl,
			},
		},
	},
	{
		ID:       "chuan@yu.com",
		Password: "n7L2LqtW3kpv#RF%",
		Req: Request{
			URL:    v2freeLoginUrl,
			Method: http.MethodPost,
			Headers: map[string]string{
				"content-type": "application/x-www-form-urlencoded; charset=UTF-8",
				"referer":      v2freeLoginUrl,
			},
		},
	},
}

func (acc *Account) getCookie(cookies []string) string {
	var sb strings.Builder
	for _, v := range cookies {
		parts := strings.Split(v, ";")
		uid := strings.TrimSpace(parts[0])
		sb.WriteString(uid)
		sb.WriteString("; ")
	}

	str := sb.String()
	str = strings.TrimRight(str, "; ")

	return str
}

type HTTPRequester interface {
	SendRequest() (*http.Response, error)
	GetResponseBody(res *http.Response) ([]byte, error)
	GetResponseFunc() (*http.Response, []byte, error)
}

var _ HTTPRequester // 确保实现了 HTTPRequester接口

type Request struct {
	URL     string
	Method  string
	Payload io.Reader
	Headers map[string]string
}

func MapToJson(param map[string]string) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func (req *Request) sendRequest() (*http.Response, error) {
	client := &http.Client{}

	httpReq, err := http.NewRequest(req.Method, req.URL, req.Payload)
	if err != nil {
		return nil, err
	}

	for key, value := range req.Headers {
		httpReq.Header.Add(key, value)
	}

	res, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (req *Request) SendRequest() (*http.Response, error) {
	var res *http.Response
	var err error

	// 设置请求超时时间为 5 秒钟
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用 for 循环进行超时重试
	for i := 0; i < 3; i++ {
		// 发送请求
		res, err = req.sendRequest()
		if err == nil {
			break
		}

		// 如果请求失败，则等待一段时间后再次尝试
		fmt.Printf("Request failed on attempt %d: %v\n", i+1, err)
		select {
		case <-ctx.Done():
			// 超时或主动取消
			return nil, fmt.Errorf("request timed out: %v", ctx.Err())
		case <-time.After(2 * time.Second):
			// 等待 2 秒钟后再次尝试
		}
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (req *Request) GetResponseBody(res *http.Response) ([]byte, error) {
	body, err := io.ReadAll(res.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (req *Request) GetAccountResponseFunc(acc *Account) (*http.Response, []byte, error) {
	req.Headers["cookie"] = acc.Cookie // 装饰器模式，对 GetResponseFunc 进行装饰增加设置cookie的设置以重用HTTPRequest接口

	return req.GetResponseFunc()
}

func (req *Request) GetResponseFunc() (*http.Response, []byte, error) {
	res, err := req.SendRequest()
	if err != nil {
		return nil, nil, err
	}
	body, err := req.GetResponseBody(res)
	if err != nil {
		return nil, nil, err
	}
	return res, body, nil
}

type pushMap map[string]string

func (push pushMap) CheckinGla() {
	wg := sync.WaitGroup{}

	for _, account := range gladosAccounts {
		wg.Add(1)
		_, body, err := account.Req.GetAccountResponseFunc(&account)
		if err != nil {
			push[account.ID] = err.Error()
			continue
		}

		result := Result{}
		_ = json.Unmarshal(body, &result)
		fmt.Println(account.ID, ":", result.Message)
		push[account.ID] = result.Message
		wg.Done()
	}
	wg.Wait()
}

func (push pushMap) CheckinV2f() {
	wg := sync.WaitGroup{}

	for _, account := range v2freeAccounts {
		wg.Add(1)
		go func(acc Account) {
			acc.Req.Payload = strings.NewReader(fmt.Sprintf("email=%s&passwd=%s&code=", url.QueryEscape(acc.ID), url.QueryEscape(acc.Password)))

			res, err := acc.Req.SendRequest()
			if err != nil {
				push[acc.ID] = err.Error()
				return
			}

			acc.Req.URL = v2freeCheckUrl
			acc.Req.Payload = strings.NewReader("")
			acc.Req.Headers["referer"] = v2freeUserUrl
			cookies := res.Header["Set-Cookie"]
			acc.Cookie = acc.getCookie(cookies)

			_, body, err := acc.Req.GetAccountResponseFunc(&acc)

			if err != nil {
				push[acc.ID] = err.Error()
				return
			}

			result := Result{}
			err = json.Unmarshal(body, &result)
			if err != nil {
				push[acc.ID] = err.Error()
				return
			}
			fmt.Println(acc.ID, ":", result.Msg)
			push[acc.ID] = result.Msg
			wg.Done()
		}(account)
	}
	wg.Wait()
}

func main() {
	push := pushMap{}
	push.CheckinGla()
	push.CheckinV2f()

	req := Request{
		URL:     fmt.Sprintf("http://www.pushplus.plus/send?token=%s&content=%s&title=GlaDOS%s", pushToken, url.QueryEscape(MapToJson(push)), url.QueryEscape("GlaDOS签到")),
		Method:  http.MethodGet,
		Payload: nil,
		Headers: nil,
	}

	_, _ = req.SendRequest()
}
