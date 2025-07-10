package utils

import (
	"github.com/valyala/fasthttp"
)

func Request(method string, json []byte, reqURL string) ([]byte, int) {
	client := &fasthttp.Client{}
	if method == "" {
		method = "GET"
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(reqURL)
	req.Header.SetMethod(method)

	req.Header.Set("Content-Type", "application/json")

	if json != nil {
		req.SetBody(json)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := client.Do(req, resp); err != nil {
		return resp.Body(), resp.StatusCode()
	}

	return resp.Body(), resp.StatusCode()
}