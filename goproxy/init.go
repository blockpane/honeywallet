package goproxy

import (
	"fmt"
	"net/http"
	"time"
)

var Client = hClient()

// setup a connection pool with high number of idle connections:
func hClient() *http.Client {
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = 100
	defaultTransport.MaxIdleConnsPerHost = 100
	return &http.Client{
		Timeout: time.Second * 10,
		Transport: &defaultTransport,
	}
}
