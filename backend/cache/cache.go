package cache

import (
	"container/list"
	"sync"
	"time"
)

type cacheItem struct {
    key        string
    value      string
    expiration time.Time
}

type LRUCache struct {
    capacity int
    items    map[string]*list.Element
    order    *list.List
    mutex    sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity: capacity,
        items:    make(map[string]*list.Element),
        order:    list.New(),
    }
}

func (c *LRUCache) Get(key string) (string, bool) {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if element, found := c.items[key]; found {
        item := element.Value.(*cacheItem)
        if item.expiration.After(time.Now()) {
            c.order.MoveToFront(element)
            return item.value, true
        }
        c.order.Remove(element)
        delete(c.items, key)
    }
    return "", false
}

func (c *LRUCache) Set(key string, value string, ttl time.Duration) {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    if element, found := c.items[key]; found {
        c.order.MoveToFront(element)
        element.Value.(*cacheItem).value = value
        element.Value.(*cacheItem).expiration = time.Now().Add(ttl)
    } else {
        if c.order.Len() >= c.capacity {
            oldest := c.order.Back()
            if oldest != nil {
                c.order.Remove(oldest)
                delete(c.items, oldest.Value.(*cacheItem).key)
            }
        }
        item := &cacheItem{key: key, value: value, expiration: time.Now().Add(ttl)}
        c.items[key] = c.order.PushFront(item)
    }
}
