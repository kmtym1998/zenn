package service

import (
	"net/http"
)

type ReqHeaders struct {
	Key   string
	Value string
}

func SendRequest(method string, url string, headers []ReqHeaders) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for _, rh := range headers {
		req.Header.Set(rh.Key, rh.Value)
	}

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
