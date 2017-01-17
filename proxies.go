package main

import (
	"net/http/httputil"
	"net/url"
)

var (
	snfusionProxy *httputil.ReverseProxy
)

func init() {
	snfusionProxy = httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "clrbinetsrv.in2p3.fr:7071",
	})
}
