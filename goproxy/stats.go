package goproxy

import (
	"encoding/json"
	"fmt"
	"github.com/montanaflynn/stats"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"time"
)

func StatsHandler(l chan LogEntry) {
	db, err := dbInit()
	if err != nil {
		log.Fatal("cannot init db", err)
	}
	defer db.Close()
	f, err := os.OpenFile(`/opt/goproxy/logs/attacks.json`, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("can't open log file, %v\n", err)
	}
	defer f.Close()
	nl := byte(0x0a)
	go statsWorker(db)
	for record := range l {
		record.Request = make(map[string]interface{})
		err = json.Unmarshal(record.ReqBytes, &record.Request)
		if err != nil {
			log.Println(`error processing raw request into map`, err)
			continue
		}
		// write log
		full, e := record.Marshall()
		if e != nil {
			log.Printf("problem marshalling log entry, %v", e)
			continue
		}
		_, e = f.Write(append(full, nl))
		if e != nil {
			log.Printf("problem writing log entry, %v\n", e)
		}
		results := &resultsStub{
			Params: []*paramsStub{},
		}
		j, _ := record.MarshallRequest()
		_ = json.Unmarshal(j, results)
		// Process IP addr
		ipRecord := &IpAddress{}
		j, err := dbFetch(`ip`, []byte(record.Ip), db)
		if err != nil {
			log.Println("error getting ip from bolt,", err)
		}
		_ = json.Unmarshal(j, ipRecord)
		if len(results.Params) > 0 && len(results.Params[0].To) > 0 {
			if ipRecord.Accounts == nil {
				ipRecord.Accounts = make(map[string]int)
			}
			ipRecord.Accounts[results.Params[0].To] = ipRecord.Accounts[results.Params[0].To] + 1
		}
		if len(results.Method) > 0 {
			if ipRecord.Methods == nil {
				ipRecord.Methods = make(map[string]int)
			}
			ipRecord.Methods[results.Method] = ipRecord.Methods[results.Method] + 1
		}
		ipRecord.LastSeen = time.Now().UTC().String()
		j, _ = json.Marshal(ipRecord)
		_ = dbWriter(`ip`, []byte(record.Ip), j, db)
		// Process Address
		if len(results.Params) > 0 && len(results.Params[0].To) > 0 {
			j, err = dbFetch(`address`, []byte(results.Params[0].To), db)
			if err == nil {
				addr := &Address{}
				_ = json.Unmarshal(j, addr)
				if addr.Ips == nil {
					addr.Ips = make(map[string]int)
				}
				addr.Ips[record.Ip] = addr.Ips[record.Ip] + 1
				if addr.Methods == nil {
					addr.Methods = make(map[string]int)
				}
				addr.Methods[results.Method] = addr.Methods[results.Method] + 1
				addr.LastSeen = time.Now().UTC().String()
				j, _ = json.Marshal(addr)
				_ = dbWriter(`address`, []byte(results.Params[0].To), j, db)
			} else {
				log.Println(err)
			}
		}
		// Process method
		if len(results.Method) > 0 {
			j, err = dbFetch(`method`, []byte(results.Method), db)
			if err == nil {
				method := &Method{}
				_ = json.Unmarshal(j, method)
				if method.Ips == nil {
					method.Ips = make(map[string]int)
				}
				method.Ips[record.Ip] = method.Ips[record.Ip] + 1
				if len(results.Params) > 0 && len(results.Params[0].To) > 0 {
					if method.Accounts == nil {
						method.Accounts = make(map[string]int)
					}
					method.Accounts[results.Params[0].To] = method.Accounts[results.Params[0].To] + 1
				}
				method.LastSeen = time.Now().UTC().String()
				j, _ = json.Marshal(method)
				_ = dbWriter(`method`, []byte(results.Method), j, db)
			} else {
				log.Println(err)
			}
		}
	}
}

func dbInit() (*bolt.DB, error) {
	db, err := bolt.Open("/opt/goproxy/logs/logs.db", 0666, nil)
	if err != nil {
		return nil, err
	}
	for _, bucket := range []string{
		`ip`,
		`method`,
		`address`,
	} {
		e := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucket([]byte(bucket))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})
		if e != nil {
			log.Println("problem creating bolt buckets,", e)
		}
	}
	return db, nil
}

func dbFetch(bucket string, key []byte, db *bolt.DB) ([]byte, error) {
	var v []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v = b.Get(key)
		return nil
	})
	return v, err
}

func dbWriter(bucket string, key []byte, value []byte, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(key, value)
		return err
	})
}

func DumpStats(db *bolt.DB) (body []byte) {
	// ips
	ipBytes := make(map[string][]byte)
	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ip"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			ipBytes[string(k)] = v
		}
		return nil
	})
	ips := make(map[string]*IpAddress)
	for k, v := range ipBytes {
		i := &IpAddress{}
		_ = json.Unmarshal(v, i)
		ips[k] = i
	}
	// methods
	methodBytes := make(map[string][]byte)
	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("method"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			methodBytes[string(k)] = v
		}
		return nil
	})
	methods := make(map[string]*Method)
	for k, v := range methodBytes {
		m := &Method{}
		_ = json.Unmarshal(v, m)
		methods[k] = m
	}
	// addresses
	addressBytes := make(map[string][]byte)
	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("address"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			addressBytes[string(k)] = v
		}
		return nil
	})
	addresses := make(map[string]*Address)
	for k, v := range addressBytes {
		a := &Address{}
		_ = json.Unmarshal(v, a)
		addresses[k] = a
	}
	j, _ := json.MarshalIndent(Stats{
		Ips:       ips,
		Methods:   methods,
		Addresses: addresses,
	}, "", "  ")
	return j
}

func DumpGraph(db *bolt.DB) (dotGraph string) {
	dotGraph = "graph IP {\nlayers=\"ip:addr:link\";\n"
	ipBytes := make(map[string][]byte)
	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ip"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			ipBytes[string(k)] = v
		}
		return nil
	})
	addrBytes := make(map[string][]byte)
	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("address"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			addrBytes[string(k)] = v
		}
		return nil
	})
	addrs := make(map[string]*Address)
	for k, v := range addrBytes {
		a := &Address{}
		_ = json.Unmarshal(v, a)
		addrs[k] = a
	}
	// traverse twice, first figure out how prevalent an IP is, and if there are multiple links to an address:
	hits := make([]float64, 0)
	addrLinkCount := make(map[string]int)
	ipCount := make(map[string]int)
	ipRank := make(map[string]int)
	// first find ninety, and 99th percentile:
	for k, v := range addrs {
		if len(v.Ips) == 0 {
			continue
		}
		for ip := range v.Ips {
			ipCount[ip] = ipCount[ip] + 1
			ipRank[ip] = ipRank[ip] + v.Ips[ip]
			addrLinkCount[k] = addrLinkCount[k] + 1
		}
		for _, c := range ipCount {
			hits = append(hits, float64(c))
		}
	}
	ninety, _ := stats.Percentile(hits, 90.0)
	top, _ := stats.Percentile(hits, 99.0)
	for k, v := range addrs {
		if len(v.Ips) == 0 {
			continue
		}
		// build our IP nodes:
		for ip := range ipCount {
			ipColor := `#bababa`
			switch {
			case float64(ipRank[ip]) > top:
				ipColor = `#bf1c1c,fontcolor=white`
			case float64(ipRank[ip]) > ninety:
				ipColor = `#d3931b`
			}
			dotGraph = dotGraph + fmt.Sprintf("    \"%s\" [layer=ip,weight=%d,value=%d,color=%s];\n", ip, ipRank[ip] ^ 2, ipRank[ip], ipColor)
		}
		color := `"#494949",fontcolor="#afbbbf"`
		switch {
		case addrLinkCount[k] > 2:
			color = `"#cc1010",fontcolor=white`
		case addrLinkCount[k] > 1:
			color = `"#efa51a"`
		}
		dotGraph = dotGraph + fmt.Sprintf("    \"%s\" [layer=addr,shape=box,fillcolor=%s];\n", k, color)
		// now traverse again for links
		for ip := range v.Ips {
			ipColor := `#bababa`
			switch {
			case addrLinkCount[k] > 3:
				ipColor = `#7f1000`
			case addrLinkCount[k] > 1:
				ipColor = `#995600`
			}
			dotGraph = dotGraph + fmt.Sprintf("    \"%s\" -- %s [layer=link,value=%d,color=%s];\n", k, ip, addrLinkCount[k] / 2, ipColor)
		}
	}
	dotGraph = dotGraph + "}\n"
	return
}


func statsWorker(db *bolt.DB) {
	for {
		g, err := os.Create(`/opt/goproxy/logs/ip-account.gv`)
		if err != nil {
			log.Println(err)
		}
		f, err := os.Create(`/opt/goproxy/logs/stats.json`)
		if err != nil {
			log.Println(err)
		}
		_, _ = g.WriteString(DumpGraph(db))
		_ = g.Close()
		_, _ = f.WriteString(string(DumpStats(db)))
		_ = f.Close()
		time.Sleep(time.Minute * 2)
	}
}
