package rpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var UserAgent = "Golang qiniu/rpc package"

// --------------------------------------------------------------------

type Client struct {
	*http.Client
}

var DefaultClient = Client{http.DefaultClient}

// --------------------------------------------------------------------

type Logger interface {
	ReqId() string
	Xput(logs []string)
}

// --------------------------------------------------------------------

func (r Client) Get(l Logger, url string) (resp *http.Response, err error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	return r.Do(l, req)
}

func (r Client) PostWith(
	l Logger, url1 string, bodyType string, body io.Reader, bodyLength int) (resp *http.Response, err error) {

	req, err := http.NewRequest("POST", url1, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)
	req.ContentLength = int64(bodyLength)
	return r.Do(l, req)
}

func (r Client) PostWith64(
	l Logger, url1 string, bodyType string, body io.Reader, bodyLength int64) (resp *http.Response, err error) {

	req, err := http.NewRequest("POST", url1, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyType)
	req.ContentLength = bodyLength
	return r.Do(l, req)
}

func (r Client) PostWithForm(
	l Logger, url1 string, data map[string][]string) (resp *http.Response, err error) {

	msg := url.Values(data).Encode()
	return r.PostWith(l, url1, "application/x-www-form-urlencoded", strings.NewReader(msg), len(msg))
}

func (r Client) PostWithJson(
	l Logger, url1 string, data interface{}) (resp *http.Response, err error) {

	msg, err := json.Marshal(data)
	if err != nil {
		return
	}
	return r.PostWith(l, url1, "application/json", bytes.NewReader(msg), len(msg))
}

func (r Client) Do(l Logger, req *http.Request) (resp *http.Response, err error) {

	if l != nil {
		req.Header.Set("X-Reqid", l.ReqId())
	}

	req.Header.Set("User-Agent", UserAgent)
	resp, err = r.Client.Do(req)
	if err != nil {
		return
	}

	if l != nil {
		details := resp.Header["X-Log"]
		if len(details) > 0 {
			l.Xput(details)
		}
	}
	return
}

// --------------------------------------------------------------------

type ErrorInfo struct {
	Err     string   `json:"error"`
	Reqid   string   `json:"reqid"`
	Details []string `json:"details"`
	Code    int      `json:"code"`
}

func (r *ErrorInfo) Error() string {
	msg, _ := json.Marshal(r)
	return string(msg)
}

// --------------------------------------------------------------------

type errorRet struct {
	Error string `json:"error"`
}

func ResponseError(resp *http.Response) (err error) {

	e := &ErrorInfo{
		Details: resp.Header["X-Log"],
		Reqid:   resp.Header.Get("X-Reqid"),
		Code:    resp.StatusCode,
	}
	if resp.StatusCode > 299 {
		if resp.ContentLength != 0 {
			if ct, ok := resp.Header["Content-Type"]; ok && ct[0] == "application/json" {
				var ret1 errorRet
				json.NewDecoder(resp.Body).Decode(&ret1)
				e.Err = ret1.Error
			}
		}
	}
	return e
}

func callRet(l Logger, ret interface{}, resp *http.Response) (err error) {

	defer resp.Body.Close()

	if resp.StatusCode/100 == 2 {
		if ret != nil && resp.ContentLength != 0 {
			err = json.NewDecoder(resp.Body).Decode(ret)
			if err != nil {
				return
			}
		}
		if resp.StatusCode == 200 {
			return nil
		}
	}
	return ResponseError(resp)
}

func (r Client) CallWithForm(l Logger, ret interface{}, url1 string, param map[string][]string) (err error) {

	resp, err := r.PostWithForm(l, url1, param)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

func (r Client) CallWithJson(l Logger, ret interface{}, url1 string, param interface{}) (err error) {

	resp, err := r.PostWithJson(l, url1, param)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

func (r Client) CallWith(
	l Logger, ret interface{}, url1 string, bodyType string, body io.Reader, bodyLength int) (err error) {

	resp, err := r.PostWith(l, url1, bodyType, body, bodyLength)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

func (r Client) CallWith64(
	l Logger, ret interface{}, url1 string, bodyType string, body io.Reader, bodyLength int64) (err error) {

	resp, err := r.PostWith64(l, url1, bodyType, body, bodyLength)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

func (r Client) Call(
	l Logger, ret interface{}, url1 string) (err error) {

	resp, err := r.PostWith(l, url1, "application/x-www-form-urlencoded", nil, 0)
	if err != nil {
		return err
	}
	return callRet(l, ret, resp)
}

// --------------------------------------------------------------------
