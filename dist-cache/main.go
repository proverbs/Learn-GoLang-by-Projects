package main

import (
	"fmt"
	"proverbs.top/distcache"
	"log"
	"net/http"
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

func main() {
	distcache.NewGroup("scores", 2<<10, dbGetter)

	addr := "localhost:9999"
	peers := distcache.NewHTTPPool(addr)
	log.Println("distcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
