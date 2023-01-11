package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const METHODPOST = "POST"
const METHODGET  = "GET"

type RequestConfig struct {
	Url    string
	Header map[string]string
	Method string
	Data   map[string]string
}

type HttpRequest struct {
	Body string
	Code int
}


func maptoUrlParams(Data map[string]string) string {
	var Params string
	for key, value := range Data {
		Params += key + "=" + url.QueryEscape(value) + "&"
	}
	if Params != "" {
		if strings.HasSuffix(Params, "&") {
			Params = Params[:len(Params)-len("&")]
		}
	}
	return Params
}



func MaptoUrlPostParams(Data map[string]string) string {
	var Params string
	for key, value := range Data {
		Params += key + "=" + value + "&"
	}
	if Params != "" {
		if strings.HasSuffix(Params, "&") {
			Params = Params[:len(Params)-len("&")]
		}
	}
	return Params
}

func urlAddParams(Url string, Params string) string {
	if strings.Index(Params, "?") != -1 {
		return Url + "&" + Params
	} else {
		return Url + "?" + Params
	}
}



func Send(config RequestConfig) (HttpRequest, error) {
	Url := config.Url
	Method := config.Method
	Data := config.Data

	Header := make(map[string]string)
	if len(config.Header) > 0 {
		for key, value := range config.Header {
			key = strings.ToTitle(key)
			Header[key] = value
		}
	}

	httpRequest := HttpRequest{}

	var request *http.Request
	var err error
	var body []byte
	if Method == METHODGET {
		if len(Data) > 0 {
			Url = urlAddParams(Url, maptoUrlParams(Data))
		}
		request, err = http.NewRequest(Method, Url, nil)
		if err != nil {
			return httpRequest, err
		}
	}



	if Method == METHODPOST {
		var reader io.Reader
		if len(Data) > 0 {
			value, _ := Header["CONTENT-TYPE"]
			if strings.ToLower(value) == "application/json" {
				var  bytesData []byte
				bytesData, err = json.Marshal(Data)
				if err != nil {
					return httpRequest, err
				}
				reader = bytes.NewReader(bytesData)
			} else {
				reader = strings.NewReader(MaptoUrlPostParams(Data))
			}
		}
		request, err = http.NewRequest(Method, Url, reader)
	}
	if err != nil {
		return httpRequest, err
	}

	if len(config.Header) > 0 {
		for key, value := range config.Header {
			request.Header.Set(key, value)
		}
	}

	client := &http.Client{
		Timeout:time.Duration(30) * time.Second,
	}
	var resp *http.Response
	resp, err = client.Do(request)
	defer client.CloseIdleConnections()

	if err != nil {
		if resp != nil {
			httpRequest.Body = ""
			httpRequest.Code = resp.StatusCode
		}
		return httpRequest, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	httpRequest.Body = string(body)

	if resp.StatusCode != 200 && resp.StatusCode != 301 && resp.StatusCode != 302 {
		logrus.Warn(fmt.Sprintf("地址请求失败\n请求地址:%s\n响应code:%d\n响应body:%s",Url,resp.StatusCode,httpRequest.Body))
		return httpRequest,err
	}

	defer resp.Body.Close()


	httpRequest.Code = resp.StatusCode

	return httpRequest,err
}




type VersionData struct {
	Version   string `json:"version"`
	Versionpath    string `json:"versionpath"`
}

func GetVersionData(url string) (res VersionData,err error) {
	res = VersionData{}
	header := make(map[string]string)
	header["content-type"] = "application/x-www-form-urlencoded"

	requestConfig := RequestConfig{
		Url:    url,
		Method: METHODGET,
		Header: header,
	}

	httpRequest, err := Send(requestConfig)
	if err != nil {
		return
	}

	if httpRequest.Code != 200 {
		err = errors.New(fmt.Sprintf("err:   url->%s code->%d body->%s",url,httpRequest.Code,httpRequest.Body))
		return
	}

	err = json.Unmarshal([]byte(httpRequest.Body), &res)
	if err != nil {
		return
	}

	return
}
