package goproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandleRequest(res http.ResponseWriter, req *http.Request) {
	abort := func() {
		_, _ = res.Write([]byte{})
		res.WriteHeader(200)
	}
	if req.Method != "POST" {
		abort()
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		abort()
		return
	}
	if len(body) == 0 {
		abort()
		return
	}
	ip := strings.Split(req.RemoteAddr, ":")[0]
	// block persistent attackers by appearing to be down
	if BlockPersistent(ip) {
		res.WriteHeader(500)
		return
	}
	var jp fastjson.Parser
	b, err := jp.Parse(string(body))
	if err != nil {
		log.Printf("error parsing json request, %v\n", err)
		abort()
		return
	}
	// have to unmarshal to shoehorn into the log struct
	logIt := func(val *fastjson.Value) {
		orig := val.MarshalTo(nil)
		p, _ := strconv.ParseInt(strings.Split(req.RemoteAddr, ":")[1], 10, 64)
		logStruct := LogEntry{
			Ip:       ip,
			Port:     p,
			UnixTime: time.Now().UTC().Unix(),
			ReqBytes: orig,
		}
		LogChan <- logStruct
	}
	if "array" == fmt.Sprintf("%v", b.Type()) {
		vals, _ := b.Array()
		for _, l := range vals {
			logIt(l)
		}
	} else {
		logIt(b)
	}
	if reply, found := FakeResponse(b); found {
		_, _ = res.Write(reply)
		res.WriteHeader(200)
		return
	} else {
		orig := b.MarshalTo(nil)
		gethReply, err := ProxyRequest(GethUrl, orig)
		if err != nil {
			abort()
			return
		}
		_, _ = res.Write(gethReply)
	}
	//fmt.Println(b)
}

func FakeResponse(request *fastjson.Value) (body []byte, ok bool) {
	parse := func(b *fastjson.Value, arrResp bool) (resp []byte, found bool) {
		if b.Exists(`method`) && b.Exists(`id`) {
			switch string(b.GetStringBytes(`method`)) {
			case `personal_listWallets`:
				bb := make([]byte, 0)
				var err error
				if arrResp {
					f := []FakeWallet{NewFakeWallet(b.GetInt(`id`))}
					bb, err = json.Marshal(f)
				} else {
					f := NewFakeWallet(b.GetInt(`id`))
					bb, err = json.Marshal(f)
				}
				if err == nil {
					return bb, true
				}
			case `eth_accounts`:
				bb := make([]byte, 0)
				var err error
				if arrResp {
					f := []FakeEthAccount{NewFakeAccount(b.GetInt(`id`))}
					bb, err = json.Marshal(f)
				} else {
					f := NewFakeAccount(b.GetInt(`id`))
					bb, err = json.Marshal(f)
				}
				if err == nil {
					return bb, true
				}
				// TODO: add more intercepted calls
			}
		}
		return nil, false
	}
	if "array" == fmt.Sprintf("%v", request.Type()) {
		arr, err := request.Array()
		if err != nil {
			return
		}
		for _, a := range arr {
			r, o := parse(a, true)
			if o {
				return r, o
			}
		}
	}
	return parse(request, false)
}

func ProxyRequest(geth string, body []byte) (resp []byte, err error) {
	reply, err := Client.Post(geth, `application/json`, bytes.NewReader(body))
	if err != nil {
		return
	}
	defer reply.Body.Close()
	return ioutil.ReadAll(reply.Body)
}

func UpdateBlocking(ip chan IpUpdate) {
	for i := range ip {
		switch i.Action {
		case "new":
			IPRateLimit[i.Ip] = &RateLimit{
				Counter: 1,
			}
		case "block":
			IPRateLimit[i.Ip].AllowAfter = time.Now().Unix() + int64(WaitTime)
			IPRateLimit[i.Ip].Counter = 0
		case "update":
			IPRateLimit[i.Ip].Counter = IPRateLimit[i.Ip].Counter + 1
		}
	}
}

func BlockPersistent(ip string) bool {
	switch {
	case IPRateLimit[ip] == nil:
		IPRateLimitChannel <- IpUpdate{Ip: ip, Action: "new"}
		return false
	case IPRateLimit[ip].AllowAfter > time.Now().Unix():
		return true
	case IPRateLimit[ip].Counter > MaxRequests:
		log.Println("blocking persistent attacker", ip)
		IPRateLimitChannel <- IpUpdate{Ip: ip, Action: "block"}
		return true
	}
	IPRateLimitChannel <- IpUpdate{Ip: ip, Action: "update"}
	return false
}

