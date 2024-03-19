package Cache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const DEFAULT_BASE_PATH = "/_cache/"

// HttpPool Http client pool.
type HttpPool struct {
	// host:port
	self string
	// request path
	basePath string
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: DEFAULT_BASE_PATH,
	}
}
func (p *HttpPool) Log(format string, v ...any) {
	log.Printf("[Server %s] %s\n", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP implement http.Handler.
func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// cache group.
	groupName := parts[0]
	// the key of request data.
	key := parts[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: ", http.StatusNotFound)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.Byteslice())
}
