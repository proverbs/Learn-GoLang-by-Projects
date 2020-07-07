package distcache

import "container/list"

type Cache struct {
	maxBytes int64 // maxBytes =0 --> no limit
	nbytes int64
	ll *list.List
	cache map[string]*list.Element
	onEvicted func(string, Val)
}

type entry struct {
	key string // for deletion
	val Val
}

type Val interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Val)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

func (c *Cache) Add(key string, newVal Val) {
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		ey, _ := e.Value.(*entry) // type assertion
		c.nbytes += int64(newVal.Len()) - int64(ey.val.Len())
		ey.val = newVal
	} else {
		e = c.ll.PushFront(&entry{key, newVal})
		c.cache[key] = e
		c.nbytes += int64(len(key)) + int64(newVal.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.Evict()
	}
}

func (c *Cache) Get(key string) (Val, bool) {
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		ey, _ := e.Value.(*entry)
		return ey.val, true
	}
	return nil, false
}

func (c *Cache) Evict() {
	e := c.ll.Back()
	if e != nil {
		c.ll.Remove(e)
		ey, _ := e.Value.(*entry)
		delete(c.cache, ey.key)
		c.nbytes -= int64(len(ey.key)) + int64(ey.val.Len())
		if c.onEvicted != nil {
			c.onEvicted(ey.key, ey.val)
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
