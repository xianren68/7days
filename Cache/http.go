package Cache

import (
	"Cache/consistenthash"
	"Cache/pb"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
)

const (
	DEFAULT_BASE_PATH = "/_cache/"
	// DEFAULT_REPLICAS default number of one real node to corrsponding how many virtual nodes.
	DEFAULT_REPLICAS = 50
)

type HttpGetter struct {
	baseUrl string
}

// Get get value from remotely nodes.
func (h *HttpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf("%s%s%s", h.baseUrl, url.QueryEscape(in.GetGroup()), url.QueryEscape(in.GetKey()))
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(bytes, out)
	if err != nil {
		return fmt.Errorf("Cannot unmarshal response: %v", err)
	}
	return nil
}

// HttpPool Http client pool.
type HttpPool struct {
	// host:port
	self string
	// request path
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*HttpGetter
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
	// /<basepath>/<groupname>/<key> required.
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
	body, err := proto.Marshal(&pb.Response{Value: view.Byteslice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

// Set set consistent hash map.
func (p *HttpPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.NewMap(DEFAULT_REPLICAS, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*HttpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &HttpGetter{baseUrl: peer + p.basePath}
	}
}

func (p *HttpPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	//
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}
