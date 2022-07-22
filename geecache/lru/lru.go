package lru

import "container/list"

type Cache struct {
	ll        *list.List
	maxByte   int64
	nByte     int64
	cache     map[string]*list.Element
	onEvicted func(key string, value Value)
}

func New(maxByte int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		ll:        list.New(),
		maxByte:   maxByte,
		onEvicted: onEvicted,
		cache:     make(map[string]*list.Element),
	}
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		kv := e.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)

		c.nByte -= int64(len(kv.key) + kv.value.Len())
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nByte += int64(value.Len() - kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nByte += int64(len(key) + value.Len())
	}

	for c.maxByte != 0 && c.maxByte < c.nByte {
		c.RemoveOldest()
	}
}
