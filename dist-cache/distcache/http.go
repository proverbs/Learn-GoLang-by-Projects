package distcache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath ="/_distcache/"

type HTTPPool struct {
	selfUrl string
	basePath string
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
