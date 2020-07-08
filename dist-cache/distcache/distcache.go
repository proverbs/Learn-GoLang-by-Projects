package distcache

import (
	"fmt"
	"log"
	"sync"
)

type Group struct {
	name string
	getter Getter
	cache Cache
}

type Getter interface {
	Get(key string) ([]byte, error)
}

// type callback, is a Getter
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name: name,
		getter: getter,
		cache: Cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.cache.get(key); ok {
		log.Println("[distcache] hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	// if getFromPeer(key) failed, then:
	return g.getFromLocal(key)
}

func (g *Group) getFromLocal(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	val := ByteView{data: cloneBytes(bytes)}
	g.populateCache(key,val)
	return val, nil
}

func (g *Group) populateCache(key string, val ByteView) {
	g.cache.add(key, val)
}
