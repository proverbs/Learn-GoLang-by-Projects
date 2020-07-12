package distcache

import (
	"fmt"
	"log"
	"proverbs.top/distcache/singleflight"
	"sync"
)

type Group struct {
	name string
	getter Getter
	cache Cache
	peers PeerPicker
	loader *singleflight.Group
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
		loader: &singleflight.Group{},
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

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) load(key string) (ByteView, error) {
	data, er := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if val, err := g.getFromPeer(peer, key); err == nil {
					return val, nil
				} else {
					log.Println("[GeeCache] Failed to get from peer", err)
				}
			}
		}
		// if getFromPeer(key) failed, then:
		return g.getFromLocal(key)
	})

	if er == nil {
		return data.(ByteView), nil
	}
	return ByteView{}, er
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

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{data: bytes}, nil
}

func (g *Group) populateCache(key string, val ByteView) {
	g.cache.add(key, val)
}
