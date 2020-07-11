package distcache

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"proverbs.top/distcache/consistent_hash"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_distcache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	selfUrl string
	basePath string
	mu sync.Mutex
	peers *consistent_hash.Map
	httpGetters map[string]*httpGetter
}

func NewHTTPPool(selfUrl string) *HTTPPool {
	return &HTTPPool{
		selfUrl: selfUrl,
		basePath: defaultBasePath,
	}
}

func (hp *HTTPPool) Log(format string, args ...interface{}) {
	log.Printf("[Server %s] %s", hp.selfUrl, fmt.Sprintf(format, args...))
}

func (hp *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, hp.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	hp.Log("%s %s", r.Method, r.URL.Path)

	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(hp.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName, key := parts[0], parts[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	bv, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(bv.ByteSlice())
}

func (hp *HTTPPool) Set(peers ...string) {
	hp.mu.Lock()
	defer hp.mu.Unlock()
	hp.peers = consistent_hash.New(defaultReplicas, nil)
	hp.peers.Add(peers...)
	hp.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		hp.httpGetters[peer] = &httpGetter{baseURL: peer + hp.basePath}
	}
}

// TODO: add peers, remove peers

func (hp *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	hp.mu.Lock()
	defer hp.mu.Unlock()
	if peer := hp.peers.Get(key); peer != "" && peer != hp.selfUrl {
		hp.Log("Pick peer %s", peer)
		return hp.httpGetters[peer], true
	}
	return nil, false
}

// compile time check: *HTTPPool satisfies PeerPicker interface
var _ PeerPicker = (*HTTPPool)(nil)

type httpGetter struct {
	baseURL string
}

func (hg *httpGetter) Get(group string, key string) ([]byte, error) {
	targetUrl := fmt.Sprintf(
		"%v%v/%v",
		hg.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	resp, err := http.Get(targetUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", resp.Status)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return bytes, nil
}

// compile time check: *httpGetter satisfies PeerGetter interface
var _ PeerGetter = (*httpGetter)(nil)
