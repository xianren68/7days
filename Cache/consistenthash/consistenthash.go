package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash hash function type.
type Hash func(key []byte) uint32

type Map struct {
	hash Hash
	// the number of virtual nodes corresponding to real nodes.
	replicas int
	// hash ring.
	keys []int
	// virtual node -> name of real nodes.
	hashMap map[int]string
}

// NewMap create a Map instance.
func NewMap(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		// default hash func is crc32.ChecksumIEEE.
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some real nodes to the hash.
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// each real node corresponding to virtual nodes.
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key))) // 计算hash值.
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key))) // 计算hash值.
	// binary search to find the correct index.
	idx := sort.Search(len(m.keys), func(i int) bool { return m.keys[i] >= hash })
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
