package utils

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

var httpClient = &fasthttp.Client{
    Name: "my-payment-client",
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
    MaxConnsPerHost: 1000,
}

func Request(method string, json []byte, reqURL string, ctx context.Context) ([]byte, int) {
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

	deadline, ok := ctx.Deadline()
	if !ok {
		if err := httpClient.Do(req, resp); err != nil {
			log.Printf("HTTP request (no-context) to %s failed: %v", reqURL, err)
			return nil, fasthttp.StatusServiceUnavailable
		}
	} else {
		timeout := time.Until(deadline)
		if timeout <= 0 {
			return nil, fasthttp.StatusRequestTimeout
		}

		if err := httpClient.DoTimeout(req, resp, timeout); err != nil {
			log.Printf("HTTP request to %s failed: %v", reqURL, err)
			if errors.Is(err, fasthttp.ErrTimeout) {
				return nil, fasthttp.StatusRequestTimeout
			}
			return nil, fasthttp.StatusServiceUnavailable
		}
	}

	bodyCopy := make([]byte, len(resp.Body()))
	copy(bodyCopy, resp.Body())

	return bodyCopy, resp.StatusCode()
}