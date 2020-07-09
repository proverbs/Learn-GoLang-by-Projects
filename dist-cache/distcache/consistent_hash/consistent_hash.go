package consistent_hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func(data []byte) uint32

type Map struct {
	hash HashFunc
	replicas int
	keys []int
	hashMap map[int]string // hash of virtual nodes to name of physical nodes
}

func New(replicas int, fn HashFunc) *Map {
	m := &Map{
		replicas: replicas,
		hash: fn,
		hashMap: make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash_key := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash_key)
			m.hashMap[hash_key] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) deleteKeyByIndex(idx int) {
	m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
}

func (m *Map) findKey(hash_key int) (int, bool) {
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash_key
	})

	if m.keys[idx%len(m.keys)] != hash_key {
		return 0, false
	}
	return idx, true
}

func (m *Map) Remove(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash_key := int(m.hash([]byte(strconv.Itoa(i) + key)))
			if idx, ok := m.findKey(hash_key); ok {
				m.deleteKeyByIndex(idx)
				delete(m.hashMap, hash_key)
			}
		}
	}
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
