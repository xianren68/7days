package Cache

import "Cache/pb"

// PeerPicker select nodes from hash ring.
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter get cache value by group and key.
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
