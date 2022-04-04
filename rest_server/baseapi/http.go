package baseapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/google/go-querystring/query"
)

var gClient http.Client

func MakeHttpClient(uri string, auth string, method string, body *bytes.Buffer, queryStr string) (*http.Client, *http.Request) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, nil
	}
	if len(auth) > 0 {
		req.Header.Add("Authorization", auth)
	}
	req.Header.Add("Content-Type", "application/json")
	if len(queryStr) > 0 {
		req.URL.RawQuery = queryStr
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 10000
	t.MaxConnsPerHost = 100000
	t.MaxIdleConnsPerHost = 1000
	t.DisableKeepAlives = false

	gClient.Timeout = 10 * time.Second
	gClient.Transport = t
	//client := &http.Client{
	//	Timeout:   10 * time.Second,
	//	Transport: t,
	//}

	return &gClient, req
}

func HttpCall(uri string, auth string, method string, body *bytes.Buffer, queryStruct interface{}) (*base.BaseResponse, error) {
	var v url.Values
	var queryStr string
	if queryStruct != nil {
		v, _ = query.Values(queryStruct)
		queryStr = v.Encode()
	}

	client, req := MakeHttpClient(uri, auth, method, body, queryStr)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
	}

	return ParseResponse(resp)
}

func ParseResponse(resp *http.Response) (*base.BaseResponse, error) {
	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		log.Errorf("HttpCall ParseResponse(%v)", resp.StatusCode)
		return nil, errors.New("HttpCall ParseResponse")
	}

	baseResp := new(base.BaseResponse)
	err := decoder.Decode(&baseResp)
	if err != nil {
		log.Errorf("ParseResponse Decode err: %v", err)
		return nil, err
	}
	return baseResp, err
}
