package resp3

import (
	"container/list"
	"time"
)

var timeStartingPoint int64

func init() {
	timeStartingPoint = time.Now().UnixNano()
}

// CachedValue defines the cached redis objects.
type CachedValue struct {
	Slot      uint32
	Key       string
	Value     []byte
	timestamp uint32
}

// Size returns the size of this object.
func (cv *CachedValue) Size() int {
	return len(cv.Key) + len(cv.Value) + 8
}

// TimeStamp return the timestamp of this object.
func (cv *CachedValue) TimeStamp() int64 {
	return int64(cv.timestamp)*1e6 + timeStartingPoint
}

func currentTimeStamp() uint32 {
	t := time.Now().UnixNano() - timeStartingPoint
	return uint32(t / 1e6)
}

// Cache is an LRU cache. It is not safe for concurrent access.
type Cache struct {
	// MaxMemory is the maximum memory of cache entries before
	// an item is evicted. Zero means no limit.
	MaxMemory int

	// CurrentMemory is the current used memory.
	CurrentMemory int

	// OnEvicted optionally specifies a callback function to be
	// executed when an entry is purged from the cache.
	OnEvicted func(key string, value *CachedValue)

	ll    *list.List
	cache map[interface{}]*list.Element
}

// New creates a new Cache.
// If maxEntries is zero, the cache has no limit and it's assumed
// that eviction is done by the caller.
func New(maxMemory int) *Cache {
	return &Cache{
		MaxMemory: maxMemory,
		ll:        list.New(),
		cache:     make(map[interface{}]*list.Element),
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value *CachedValue) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value = value
		return
	}

	value.timestamp = currentTimeStamp()

	ele := c.ll.PushFront(value)
	c.cache[key] = ele
	c.CurrentMemory += value.Size()
	for c.MaxMemory != 0 && c.CurrentMemory > c.MaxMemory {
		c.RemoveOldest()
	}
}

// AddValue adds a raw value to the cache.
func (c *Cache) AddValue(slot uint32, key string, value []byte) {
	cv := &CachedValue{
		Slot:  slot,
		Key:   key,
		Value: value,
	}
	c.Add(key, cv)
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key string) (value *CachedValue, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*CachedValue), true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *Cache) Remove(key string) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
		cv := ele.Value.(*CachedValue)
		c.CurrentMemory -= cv.Size()
	}
}

func (c *Cache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*CachedValue)
	delete(c.cache, kv.Key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.Key, kv)
	}
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Clear purges all stored items from the cache.
func (c *Cache) Clear() {
	if c.OnEvicted != nil {
		for _, e := range c.cache {
			kv := e.Value.(*CachedValue)
			c.OnEvicted(kv.Key, kv)
		}
	}
	c.ll = nil
	c.cache = nil
}
