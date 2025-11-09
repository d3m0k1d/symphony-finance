package httpclient

import (
	"context"
	"net/http"
)

type Client struct {
	c         http.Client
	hooksPre  []func(req *http.Request)
	hooksPost []func(req *http.Response)
}

func (c Client) Request(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	for _, f := range c.hooksPre {
		f(req)
	}
	resp, err = c.c.Do(req.WithContext(ctx))
	if err != nil {
		return
	}
	for _, f := range c.hooksPost {
		f(resp)
	}
	return
}
