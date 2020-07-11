package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"proverbs.top/distcache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

var loadCounts = make(map[string]int, len(db))

var dbGetter distcache.Getter = distcache.GetterFunc(
	func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key]++
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	})

func createGroup() *distcache.Group {
	return distcache.NewGroup("scores", 2<<10, dbGetter)
}

func startCacheServer(selfUrl string, peerAddrs []string, dcg *distcache.Group) {
	peers := distcache.NewHTTPPool(selfUrl)
	peers.Set(peerAddrs...)
	dcg.RegisterPeers(peers)
	log.Println("distcache is running at", selfUrl)
	log.Fatal(http.ListenAndServe(selfUrl[7:], peers)) // localhost:xxxx
}

func startAPIServer(apiAddr string, dcg *distcache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := dcg.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "distcache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	dcg := createGroup()
	if api {
		go startAPIServer(apiAddr, dcg)
	}
	startCacheServer(addrMap[port], addrs, dcg)
}
