package cache

import (
	"container/list"
	"sync"
)

type Storage interface {
	Set(key string, value interface{})
	Get(key string) (Record, bool)
	Remove(key string)
	Len() int
}

type Cache struct {
	mutex    *sync.Mutex
	records  map[string]*list.Element
	queue    *list.List
	capacity int
}

type Record struct {
	Key   string
	Value interface{}
}

func NewLRUCache(capacity int) *Cache {
	return &Cache{
		mutex:    new(sync.Mutex),
		capacity: capacity,
		records:  make(map[string]*list.Element, capacity),
		queue:    list.New(),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	record := Record{Key: key, Value: value}

	val, ok := c.records[key]

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if ok {
		c.queue.Remove(val)
	} else if c.queue.Len() == c.capacity {
		backKey := c.queue.Back().Value.(Record).Key
		delete(c.records, backKey)
		c.queue.Remove(c.queue.Back())
	}

	c.queue.PushFront(record)
	c.records[key] = c.queue.Front()
}

func (c *Cache) Get(key string) (Record, bool) {
	val, ok := c.records[key]

	if ok {
		record := val.Value.(Record)

		c.mutex.Lock()
		c.queue.MoveToFront(val)
		c.mutex.Unlock()

		return record, ok
	}

	return Record{}, false
}

func (c *Cache) Remove(key string) {
	val, ok := c.records[key]

	if ok {
		record := val.Value.(Record)
		c.mutex.Lock()
		delete(c.records, record.Key)
		c.queue.Remove(val)
		c.mutex.Unlock()
	}
}

func (c *Cache) Len() int {
	return c.queue.Len()
}
