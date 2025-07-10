package utils

import (
	"github.com/valyala/fasthttp"
)

func Request(method string, json []byte) ([]byte, int) {
	client := &fasthttp.Client{}
	if method == "" {
		method = "GET"
	}

	reqURL := "http://payment-processor-default:8080/payments"

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