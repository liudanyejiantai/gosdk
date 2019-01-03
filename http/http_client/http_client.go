// Copyright 2018 yejiantai Authors
//
// package http_client http客户端
package http_client

import (
	"io/ioutil"
	"net/http"
)

// http get
func HttpGet(url string) (string, error) {
	var (
		resp *http.Response
		err  error
		body []byte
	)

	if resp, err = http.Get(url); err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	}

	return string(body), nil
}
