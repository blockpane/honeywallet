package main

import (
	"gitlab.com/blockpane/honeywallet/goproxy"
	"log"
	"net/http"
	"syscall"
	"time"
)

func main() {
	http.HandleFunc("/", goproxy.HandleRequest)
	syscall.Umask(0022)
	go goproxy.StatsHandler(goproxy.LogChan)
	go goproxy.UpdateBlocking(goproxy.IPRateLimitChannel)
	go goproxy.UpdateRand()
	s := &http.Server{
		Addr:         ":8545",
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
